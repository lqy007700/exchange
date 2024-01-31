package internal

import (
	"asset-service/internal/sequencer"
	"asset-service/pkg/util"
	"asset-service/repository/mq"
	"asset-service/repository/mysql"
	"asset-service/repository/redis"
	order2 "github.com/lqy007700/exchange/common/order"

	"encoding/json"
	"go-micro.dev/v4/broker"
	"go-micro.dev/v4/logger"
	"math/big"
	"time"
)

// TrustOrder 委托单
type TrustOrder struct {
	ID       string        `json:"id"`
	UserID   int64         `json:"user_id"`
	Symbol   string        `json:"symbol"`
	Amount   *big.Float    `json:"amount"`
	Price    *big.Float    `json:"price"`
	Status   order2.Status `json:"status"`
	CreateAt time.Time     `json:"create_at"`
}

type OrderService struct {
	db            *mysql.AssetDB
	redis         *redis.AssetCache
	kafkaProducer *mq.Service
	seq           *sequencer.Seq
	mq            *mq.Service
}

func NewOrderService(db *mysql.AssetDB, redis *redis.AssetCache, mq *mq.Service) *OrderService {
	return &OrderService{db: db, redis: redis, mq: mq}
}

// CreateOrder 创建订单
func (o *OrderService) CreateOrder(order *TrustOrder) error {
	order.ID = util.GetNum()
	order.Status = order2.Pending

	producer := o.mq.Producer()

	ev := &Event{
		Type: Create,
		Data: order,
	}

	marshal, err := json.Marshal(ev)
	if err != nil {
		logger.Errorf("marshal order error: %v", err)
		return err
	}

	msg := &broker.Message{
		Header: nil,
		Body:   marshal,
	}

	err = producer.Publish("trade", msg)
	if err != nil {
		logger.Errorf("publish order error: %v", err)
		return err
	}
	return nil
}
