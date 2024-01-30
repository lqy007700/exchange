package main

import (
	"asset-service/asset-service/config"
	"asset-service/asset-service/internal"
	asset_service "asset-service/asset-service/proto"
	"asset-service/asset-service/repository/mq"
	"asset-service/asset-service/repository/mysql"
	"asset-service/asset-service/repository/redis"
	"asset-service/asset-service/rpc"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"go-micro.dev/v4"
	"go-micro.dev/v4/broker"
	"go-micro.dev/v4/logger"
	"strconv"
	"time"
)

func main() {
	err := config.Init()
	if err != nil {
		logger.Errorf("init config error: %v", err)
		return
	}

	config.InitLogger()

	// Create a new service. Optionally include some options here.
	service := micro.NewService(
		micro.Name(config.Conf.Micro.Name),
		micro.Address(config.Conf.RPCServer.Addr),
		//micro.Logger(initLogger),
	)

	db := mysql.New()
	cache := redis.NewAssetCache()

	mqSvc := mq.NewService()
	as := &rpc.AssetService{
		Asset: internal.NewAssetService(db, cache),
		Order: internal.NewOrderService(db, cache, mqSvc),
	}

	err = asset_service.RegisterAssetServiceHandler(service.Server(), as)
	if err != nil {
		logger.Errorf("register asset service error: %v", err)
		return
	}

	// Run the server
	if err := service.Run(); err != nil {
		logger.Errorf("run asset service error: %v", err)
		return
	}
}

//func main() {
//	service := mq.NewService()
//	err := service.Init()
//	if err != nil {
//		logger.Errorf("mq service init error: %v", err)
//		return
//	}
//
//	//go pub(service.Producer())
//
//	consumer := service.Consumer()
//	subOptions := []broker.SubscribeOption{broker.Queue("queue-1"), broker.DisableAutoAck()}
//	subscriber, err := consumer.Subscribe("topic0", handlerMsg, subOptions...)
//	if err != nil {
//		logger.Errorf("mq service init error: %v", err)
//		return
//	}
//
//	logger.Info(subscriber)
//	defer service.Close()
//	select {}
//}

func pub(pro broker.Broker) {
	i := 0
	for {
		time.Sleep(time.Second)
		err := pro.Publish("topic0", &broker.Message{Header: map[string]string{"id": strconv.Itoa(i)}, Body: []byte("hello world")})
		if err != nil {
			logger.Errorf("mq service init error: %v", err)
			return
		}
		i++
	}
}

func handlerMsg(event broker.Event) error {
	fmt.Println("handler msg")
	fmt.Println(event.Message().Header)
	fmt.Println(string(event.Message().Body))
	err := event.Ack()
	if err != nil {
		return err
	}
	return nil
}
