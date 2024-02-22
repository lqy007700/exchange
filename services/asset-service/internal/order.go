package internal

import (
	"asset-service/internal/sequencer"
	"asset-service/pkg/util"
	"asset-service/repository/mq"
	"asset-service/repository/mysql"
	"asset-service/repository/redis"
	order2 "github.com/lqy007700/exchange/common/order"

	"encoding/json"
	"go-micro.dev/v4/logger"
)

type OrderService struct {
	db    *mysql.AssetDB
	redis *redis.AssetCache
	seq   *sequencer.Seq
	mq    *mq.KafkaClient
}

func NewOrderService(db *mysql.AssetDB, redis *redis.AssetCache, mq *mq.KafkaClient) *OrderService {
	return &OrderService{db: db, redis: redis, mq: mq}
}

// CreateOrder 创建订单
func (o *OrderService) CreateOrder(order *order2.OrderEntity) error {
	order.ID = util.GetNum()
	order.Status = order2.Pending

	ev := &Event{
		Type: Create,
		Data: order,
	}

	marshal, err := json.Marshal(ev)
	if err != nil {
		logger.Errorf("marshal order error: %v", err)
		return err
	}

	err = o.mq.Produce("queue-engine-topic-btc_usdt1", marshal)
	if err != nil {
		logger.Errorf("publish order error: %v", err)
		return err
	}
	return nil
}
