package internal

import (
	"container/heap"
	"encoding/json"
	"engine-service/repository/mq"
	"engine-service/repository/redis"
	"fmt"
	"github.com/Shopify/sarama"
	order2 "github.com/lqy007700/exchange/common/order"
	"github.com/pkg/errors"
	"go-micro.dev/v4/logger"
	"math/big"
)

const (
	// QueueEngineTopic 撮合引擎消息队列
	QueueEngineTopic = "queue-engine-topic-%s" // coin_pair
)

type Engine struct {
	//mq   *mq.KafkaClient
	mq *mq.KafkaClient

	buy  *BuyBook
	sell *SellBook

	// Close 关闭通道
	Close <-chan struct{}
}

func NewEngine(cache *redis.BooksCache, coinPair string, mq *mq.KafkaClient) *Engine {
	// 初始化从Cache中加载订单簿
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

	buy := &BuyBook{Common: Common{data: buyBooks}}
	sell := &SellBook{Common{data: sellBooks}}
	heap.Init(buy)
	heap.Init(sell)

	return &Engine{
		buy:  buy,
		sell: sell,
		mq:   mq,
	}
}

func (e *Engine) Start(coinPair string) {
	// todo 需要处理 engine 的 close 信号
	topic := fmt.Sprintf(QueueEngineTopic, coinPair)
	e.mq.Consume(topic, e.processMsg)
	defer e.mq.Close()
}

// ProcessOrder 撮合
func (e *Engine) processOrder(takerOrder *order2.OrderEntity, makerBooks heap.Interface, anotherBooks heap.Interface) (*MatchResult, error) {
	takerUnfilledQuantity := takerOrder.Quantity
	matchRes := newMatchResult()

	for makerBooks.Len() > 0 {
		pop := heap.Pop(makerBooks)
		if pop == nil {
			// 对手盘不存在
			logger.Info("takerBooks is nil")
			break
		}

		makerOrder, ok := pop.(*order2.OrderEntity)
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
		heap.Push(anotherBooks, takerOrder)
	}
	return matchRes, nil
}

// processMsg 接收消息队列消息
func (e *Engine) processMsg(msg *sarama.ConsumerMessage) error {
	order := &order2.OrderEntity{}
	err := json.Unmarshal(msg.Value, order)
	if err != nil {
		logger.Errorf("unmarshal order error: %v", err)
		return err
	}

	if order == nil {
		logger.Error("order is nil")
		return errors.New("order is nil")
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

	mr.mq = e.mq
	mr.sendMatchResToQueue()
	return nil
}

func minFloat(x, y *big.Float) *big.Float {
	result := x.Cmp(y)
	if result == -1 {
		return x
	}
	return y
}
