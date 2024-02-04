package internal

import (
	"container/heap"
	"encoding/json"
	"engine-service/repository/mq"
	"github.com/Shopify/sarama"
	order2 "github.com/lqy007700/exchange/common/order"
	"github.com/pkg/errors"
	"go-micro.dev/v4/logger"
	"math/big"
)

// TrustOrder 委托单
//type TrustOrder struct {
//	ID               string           `json:"id"`
//	UserID           int64            `json:"user_id"`
//	Symbol           string           `json:"symbol"`
//	Amount           *big.Float       `json:"amount"`
//	Price            *big.Float       `json:"price"`
//	UnfilledQuantity *big.Float       `json:"unfilled_quantity"`
//	Status           order2.Status    `json:"status"`
//	CreateAt         time.Time        `json:"create_at"`
//	Direction        order2.Direction `json:"direction"`
//}

type EngineService struct {
	mq   *mq.KafkaClient
	buy  *BuyBook
	sell *SellBook
}

func NewEngineService(mq *mq.KafkaClient) *EngineService {
	buy := &BuyBook{}
	sell := &SellBook{}
	heap.Init(buy)
	heap.Init(sell)

	return &EngineService{mq: mq, buy: buy, sell: sell}
}

func (e *EngineService) processOrder(takerOrder *TrustOrder, makerBooks heap.Interface, anotherBooks heap.Interface) (*MatchResult, error) {
	takerUnfilledQuantity := takerOrder.Amount
	matchRes := newMatchResult()

	for makerBooks.Len() > 0 {
		pop := heap.Pop(makerBooks)
		if pop == nil {
			// 对手盘不存在
			logger.Info("takerBooks is nil")
			break
		}

		makerOrder, ok := pop.(*TrustOrder)
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

// ProcessMsg 处理创建订单消息
func (e *EngineService) ProcessMsg() {
	e.mq.Consume("order", e.processMsg)
}

func (e *EngineService) processMsg(msg *sarama.ConsumerMessage) error {
	order := &TrustOrder{}
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

	// 批量发送撮合结果给 kafka
	if mr != nil {
		for _, detail := range mr.matchDetails {
			detailJson, err := json.Marshal(detail)
			if err != nil {
				logger.Errorf("marshal match result detail error: %v", err)
				continue
			}
			err = e.mq.Produce("match_result", detailJson)
			if err != nil {
				logger.Errorf("send match result error: %v", err)
				continue
			}
		}
	}
	return nil
}

func minFloat(x, y *big.Float) *big.Float {
	result := x.Cmp(y)
	if result == -1 {
		return x
	}
	return y
}
