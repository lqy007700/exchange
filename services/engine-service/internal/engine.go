package internal

import (
	"container/heap"
	"context"
	"encoding/json"
	"engine-service/config"
	"engine-service/repository/redis"
	"fmt"
	"github.com/lqy007700/exchange/common/engine"
	"github.com/lqy007700/exchange/common/mq"
	order2 "github.com/lqy007700/exchange/common/order"
	"go-micro.dev/v4/logger"
	"math/big"
	"time"
)

type Engine struct {
	CoinPair string

	mq *mq.Consumer

	buy  *BuyBook
	sell *SellBook

	cache *redis.BooksCache

	cancel context.CancelFunc
}

func NewEngine(cache *redis.BooksCache, coinPair string) *Engine {
	// 初始化从Cache中加载订单簿
	// todo 需要检测 cache 中的订单是否正确,和 db 中的订单是否一致
	buyBooks, err := cache.GetBooks(coinPair, order2.Buy)
	if err != nil {
		logger.Errorf("get buy books error: %v", err)
		panic(err)
	}
	sellBooks, err := cache.GetBooks(coinPair, order2.Sell)
	if err != nil {
		logger.Errorf("get sell books error: %v", err)
		panic(err)
	}

	buy := &BuyBook{Common{data: buyBooks}}
	sell := &SellBook{Common{data: sellBooks}}
	heap.Init(buy)
	heap.Init(sell)

	return &Engine{
		buy:      buy,
		sell:     sell,
		CoinPair: coinPair,
		cache:    cache,
	}
}

func (e *Engine) GetOrderBookList() string {
	res := make(map[string][]*order2.OrderEntity, 2)

	res["buy"] = e.buy.data
	res["sell"] = e.sell.data

	marshal, err := json.Marshal(res)
	if err != nil {
		return ""
	}
	return string(marshal)
}

func (e *Engine) Start(coinPair string) {
	topic := fmt.Sprintf(engine.QueueEngineTopic, coinPair)

	consumer, err := mq.NewConsumer(config.Conf.Kafka.Brokers, coinPair)
	if err != nil {
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	e.cancel = cancel
	ch := &mq.ConsumerHandler{
		HandleMessage: e.processMsg,
		Ctx:           ctx,
	}

	err = consumer.Consume(ctx, []string{topic}, ch)
	if err != nil {
		return
	}
}

// ProcessOrder 撮合
func (e *Engine) processOrder(takerOrder *order2.OrderEntity, makerBooks WrapHeap, anotherBooks WrapHeap) (*MatchResult, error) {
	logger.Infof("process order: %v", takerOrder)
	takerUnfilledQuantity := takerOrder.Quantity
	matchRes := newMatchResult()

	for makerBooks.Len() > 0 {
		makerOrder, ok := makerBooks.Peek().(*order2.OrderEntity)
		if !ok || makerOrder == nil {
			logger.Error("pop is not a TrustOrder")
			break
		}

		if takerOrder.Direction == order2.Buy && takerOrder.Price.Cmp(makerOrder.Price) < 0 {
			logger.Infof("takerOrder.Price: %v, makeOrder.Price: %v", takerOrder.Price, makerOrder.Price)
			// 买入价格比卖盘第一档低
			break
		} else if takerOrder.Direction == order2.Sell && takerOrder.Price.Cmp(makerOrder.Price) > 0 {
			logger.Infof("takerOrder.Price: %v, makeOrder.Price: %v", takerOrder.Price, makerOrder.Price)
			// 卖出价格比买盘第一档高
			break
		}

		// 成交数量为两者的小值
		matchedQuantity := minFloat(takerUnfilledQuantity, makerOrder.UnfilledQuantity)

		// 成交记录
		matchRes.add(makerOrder.Price, matchedQuantity, takerOrder, makerOrder)

		// 更新成交后的订单数量:
		takerUnfilledQuantity = new(big.Float).Sub(takerUnfilledQuantity, matchedQuantity)
		makerUnfilledQuantity := new(big.Float).Sub(makerOrder.UnfilledQuantity, matchedQuantity)

		// 对手盘部分成交: 更改 make
		if makerUnfilledQuantity.Sign() > 0 {
			makerOrder.UnfilledQuantity = makerUnfilledQuantity
			heap.Push(makerBooks, makerOrder)
		} else {
			makerOrder.UnfilledQuantity = makerUnfilledQuantity
			makerOrder.Status = order2.FullyFilled
			makerBooks.Pop()
		}

		if takerUnfilledQuantity.Sign() == 0 {
			takerOrder.UnfilledQuantity = takerUnfilledQuantity
			takerOrder.Status = order2.FullyFilled
			break
		}
	}

	// Taker订单未完全成交时，放入订单簿:
	if takerUnfilledQuantity.Sign() > 0 {
		takerOrder.UnfilledQuantity = takerUnfilledQuantity

		// 区分是部分成交还是完全未成交
		status := order2.Pending
		if takerUnfilledQuantity.Cmp(takerOrder.UnfilledQuantity) == 0 {
			status = order2.PartialFilled
		}
		takerOrder.Status = status
		logger.Infof("Taker订单未完全成交: %v", takerOrder)
		heap.Push(anotherBooks, takerOrder)
	}
	return matchRes, nil
}

// processMsg 接收消息队列消息
func (e *Engine) processMsg(msg []byte) {
	ev := &Event{
		Data: &order2.OrderEntity{},
	}
	err := json.Unmarshal(msg, ev)
	if err != nil {
		logger.Errorf("unmarshal order error: %v", err)
		return
	}

	order, ok := ev.Data.(*order2.OrderEntity)
	if !ok {
		logger.Errorf("ev.Data.(*order2.OrderEntity) error: %v", err)
		return
	}

	if order == nil {
		logger.Error("order is nil")
		return
	}

	var mr *MatchResult
	switch order.Direction {
	case order2.Buy:
		mr, err = e.processOrder(order, e.sell, e.buy)
	case order2.Sell:
		mr, err = e.processOrder(order, e.buy, e.sell)
	default:
		logger.Errorf("unknown event type: %v", order.Direction)
	}

	// TODO 撮合结果
	// 1/ 发送消息等资产服务结算
	// 2/ 发送交易结果给客户端
	// 3/ 发送交易结果给行情服务
	fmt.Println(mr)

	//mr.mq = e.mq
	//mr.sendMatchResToQueue()
	return
}

// Shutdown 关闭
func (e *Engine) Shutdown() {
	if e.cancel != nil {
		logger.Info("engine is shutting down")
		e.cancel()
	}

	time.Sleep(1 * time.Second)
	_ = e.cache.SetBooks(e.CoinPair, order2.Buy, e.buy.data)
	_ = e.cache.SetBooks(e.CoinPair, order2.Sell, e.sell.data)
}

func minFloat(x, y *big.Float) *big.Float {
	result := x.Cmp(y)
	if result == -1 {
		return x
	}
	return y
}
