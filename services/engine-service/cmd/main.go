package main

import (
	"engine-service/config"
	"engine-service/internal"
	"engine-service/repository/mq"
	"engine-service/repository/redis"
	"go-micro.dev/v4"
	"go-micro.dev/v4/logger"
	"time"
)

func main() {
	err := config.Init()
	if err != nil {
		logger.Fatalf("init config error: %v", err)
		return
	}

	service := micro.NewService(
		micro.Name(config.Conf.Micro.Name),
		micro.Address(config.Conf.RPCServer.Addr),
	)
	service.Init()

	client, err := mq.NewKafkaClient()
	if err != nil {
		logger.Fatal(err)
		return
	}

	cache := redis.NewBooksCache()

	initEngine(client, cache)

	if err := service.Run(); err != nil {
		panic(err)
	}
}

// initEngine 启动撮合
func initEngine(client *mq.KafkaClient, cache *redis.BooksCache) {
	logger.Info("init engine")
	engines := make(map[string]*internal.Engine)
	// todo 交易对处理
	coinParis := []string{"btc_usdt", "eth_usdt", "eth_btc"}
	for _, coinPair := range coinParis {
		engine := internal.NewEngine(cache, coinPair, client)
		go engine.Start(coinPair)
		go printOrderBook(engine)
		engines[coinPair] = engine
		logger.Infof("init engine for %s success", coinPair)
	}
}

func printOrderBook(engine *internal.Engine) {
	ticker := time.NewTicker(20 * time.Second)
	for {
		<-ticker.C
		engine.GetOrderBookList()
	}
}
