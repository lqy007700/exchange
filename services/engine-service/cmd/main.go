package main

import (
	"engine-service/config"
	"engine-service/internal"
	"engine-service/repository/mq"
	"go-micro.dev/v4"
	"go-micro.dev/v4/logger"
)

func main() {
	err := config.Init()
	if err != nil {
		logger.Errorf("init config error: %v", err)
		return
	}

	// Create a new service. Optionally include some options here.
	service := micro.NewService(
		micro.Name(config.Conf.Micro.Name),
		micro.Address(config.Conf.RPCServer.Addr),
	)
	service.Init()

	client, err := mq.NewKafkaClient()
	if err != nil {
		return
	}

	initEngine(client)

	if err := service.Run(); err != nil {
		panic(err)
	}
}

// initEngine 初始化撮合
func initEngine(client *mq.KafkaClient) {
	logger.Info("init engine")
	engines := make(map[string]*internal.Engine)
	coinParis := []string{"btc_usdt", "eth_usdt", "eth_btc"}
	for _, coinPair := range coinParis {
		engine := internal.NewEngine(nil, coinPair, client)
		go engine.Start(coinPair)
		engines[coinPair] = engine
		logger.Infof("init engine for %s success", coinPair)
	}
}
