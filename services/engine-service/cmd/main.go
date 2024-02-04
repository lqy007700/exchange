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

	engineService := internal.NewEngineService(client)

	go engineService.ProcessMsg()

	if err := service.Run(); err != nil {
		panic(err)
	}
}
