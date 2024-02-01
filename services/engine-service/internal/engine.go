package internal

import (
	"container/heap"
	"encoding/json"
	"engine-service/repository/mq"
	"github.com/lqy007700/exchange/common/order"
	"go-micro.dev/v4/broker"
	"go-micro.dev/v4/logger"
	"math/big"
	"time"
)

// TrustOrder 委托单
type TrustOrder struct {
	ID               string          `json:"id"`
	UserID           int64           `json:"user_id"`
	Symbol           string          `json:"symbol"`
	Amount           *big.Float      `json:"amount"`
	Price            *big.Float      `json:"price"`
	UnfilledQuantity *big.Float      `json:"unfilled_quantity"`
	Status           order.Status    `json:"status"`
	CreateAt         time.Time       `json:"create_at"`
	Direction        order.Direction `json:"direction"`
}

type EngineService struct {
	mq   *mq.Service
	buy  *BuyBook
	sell *SellBook
}

func (e *EngineService) Init() {
	e.buy = &BuyBook{}
	e.sell = &SellBook{}
	heap.Init(e.buy)
	heap.Init(e.sell)

	consumer := e.mq.Consumer()
	subOptions := []broker.SubscribeOption{broker.Queue("queue-1"), broker.DisableAutoAck()}
	subscriber, err := consumer.Subscribe("topic0", e.processMsg, subOptions...)
	if err != nil {
		logger.Errorf("mq service init error: %v", err)
		return
	}

	logger.Info(subscriber)
	defer e.mq.Close()
	select {}
}

func (e *EngineService) processMsg(event broker.Event) error {
	logger.Infof("handler msg %v", event.Message())

	ev := &Event{}
	err := json.Unmarshal(event.Message().Body, ev)
	if err != nil {
		logger.Errorf("unmarshal order error: %v", err)
		return err
	}

	switch ev.Type {
	case Create:
		e.createOrder(ev.Data.(*TrustOrder))
	}
	return nil
}

func (e *EngineService) createOrder(order *TrustOrder) {
	switch order.Direction {
	case Buy:
		e.processOrder(order, e.sell, e.buy)
		break
	case Sell:
		e.processOrder(order, e.buy, e.sell)
		break
	default:
		return
	}
}

func (e *EngineService) processOrder(takerOrder *TrustOrder, makerBooks Books, anotherBooks Books) {
	takerUnfilledQuantity := takerOrder.Amount

	makerBook, ok := makerBooks.(heap.Interface)
	if !ok || makerBook == nil {
		logger.Error("makeBook is not a heap.Interface")
		return
	}

	for makerBook.Len() > 0 {
		pop := heap.Pop(makerBook)
		if pop == nil {
			// 对手盘不存在
			break
		}

		makeOrder, ok := pop.(*TrustOrder)
		if !ok || makeOrder == nil {
			logger.Error("pop is not a TrustOrder")
			break
		}

		if takerOrder.Direction == Buy && takerOrder.Price.Cmp(makeOrder.Price) < 0 {
			// 买入价格比卖盘第一档低
			break
		} else if takerOrder.Direction == Sell && takerOrder.Price.Cmp(makeOrder.Price) > 0 {
			// 卖出价格比买盘第一档高
			break
		}

		// 成交数量为两者的小值
		matchedQuantity := minFloat(takerUnfilledQuantity, makeOrder.UnfilledQuantity)
		// 成交记录 todo

		// 更新成交后的订单数量:
		takerUnfilledQuantity = new(big.Float).Sub(takerUnfilledQuantity, matchedQuantity)
		makerUnfilledQuantity := new(big.Float).Sub(makeOrder.UnfilledQuantity, matchedQuantity)

		// 对手盘完全成交后，从订单簿中删除:
		if makerUnfilledQuantity.Sign() == 0 {
			// 更改 maker 订单状态为完全成交 todo
			// 从 maker books 移除
			//makeOrder.Update(0, FullyFilled)
		} else {
			// 对手盘部分成交: 更改 make
			//r 订单状态为部分成交 todo
			//makeOrder.Update(makerUnfilledQuantity, PartialFilled)
			heap.Push(makerBook, makeOrder)
		}

		if takerUnfilledQuantity.Sign() == 0 {
			// 更改 taker 订单状态为完全成交 todo
			break
		}
	}

	// Taker订单未完全成交时，放入订单簿:
	if takerUnfilledQuantity.Sign() > 0 {
		// 更改 taker 订单状态为部分成交 和 未成交数量 todo
		anotherBook, ok := anotherBooks.(heap.Interface)
		if !ok {
			logger.Error("anotherBook is not a SellBook")
			return
		}
		if anotherBook == nil {
			// todo 初始化对手盘 books 列表
		}
		heap.Push(anotherBook, takerOrder)
	}
}

func minFloat(x, y *big.Float) *big.Float {
	result := x.Cmp(y)
	if result == -1 {
		return x
	}
	return y
}
