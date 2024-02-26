package internal

import (
	"asset-service/config"
	"asset-service/internal/sequencer"
	"asset-service/pkg/util"
	"asset-service/repository/mysql"
	"asset-service/repository/redis"
	"fmt"
	"github.com/lqy007700/exchange/common/engine"
	"github.com/lqy007700/exchange/common/mq"
	order2 "github.com/lqy007700/exchange/common/order"

	"encoding/json"
	"go-micro.dev/v4/logger"
)

type OrderService struct {
	db    *mysql.AssetDB
	redis *redis.AssetCache
	seq   *sequencer.Seq
	mq    *mq.Producer
}

func NewOrderService(db *mysql.AssetDB, redis *redis.AssetCache) *OrderService {
	producer, err := mq.NewProducer(config.Conf.Kafka.Brokers)
	if err != nil {
		logger.Fatalf("init kafka producer error: %v", err)
		return nil
	}
	return &OrderService{db: db, redis: redis, mq: producer}
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

	topic := fmt.Sprintf(engine.QueueEngineTopic, fmt.Sprintf("%s_%s", order.CoinFrom, order.CoinTo))
	o.mq.ProduceMessage(topic, string(marshal))
	return nil
}
