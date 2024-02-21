package internal

import (
	"asset-service/config"
	"asset-service/repository/mq"
	"asset-service/repository/mysql"
	"asset-service/repository/redis"
	"github.com/lqy007700/exchange/common/order"
	"go-micro.dev/v4/logger"
	"math/big"
	"testing"
	"time"
)

var orderSvc *OrderService

func TestOrderService_CreateOrder(t *testing.T) {
	err := config.Init()
	if err != nil {
		logger.Errorf("init config error: %v", err)
		return
	}
	config.InitLogger()

	db := mysql.New()
	cache := redis.NewAssetCache()

	mqSvc, _ := mq.NewKafkaClient()

	orderSvc = NewOrderService(db, cache, mqSvc)

	o := &order.OrderEntity{
		ID:               "2",
		UserID:           1,
		CoinFrom:         "btc",
		CoinTo:           "usdt",
		Direction:        order.Buy,
		Price:            big.NewFloat(10000),
		Quantity:         big.NewFloat(10),
		UnfilledQuantity: big.NewFloat(10),
		Status:           order.Pending,
		CreateAt:         time.Now(),
	}
	err = orderSvc.CreateOrder(o)
	if err != nil {
		logger.Errorf("create order error: %v", err)
		return
	}
	return
}
