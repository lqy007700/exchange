package main

import (
	"engine-service/config"
	"engine-service/internal"
	"engine-service/repository/redis"
	"go-micro.dev/v4/web"
	"net/http"

	//"github.com/asim/go-micro/v4/web"
	"go-micro.dev/v4"
	"go-micro.dev/v4/logger"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	ShutdownSignals = []os.Signal{
		os.Interrupt, os.Kill, syscall.SIGKILL, syscall.SIGSTOP,
		syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGILL, syscall.SIGTRAP,
		syscall.SIGABRT, syscall.SIGSYS, syscall.SIGTERM,
	}

	engines = make(map[string]*internal.Engine)
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

	cache := redis.NewBooksCache()
	initEngine(cache)

	webSvc := web.NewService(web.Name("my.service"), web.Address("127.0.0.1:9901"))
	webSvc.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		coinPair := r.URL.Query().Get("coin_pair")
		var list string
		if engine, ok := engines[coinPair]; ok {
			list = engine.GetOrderBookList()
		}
		w.Write([]byte(list))
	})

	err = webSvc.Start()
	if err != nil {
		panic(err)
		return
	}

	go func() {
		if err := service.Run(); err != nil {
			panic(err)
		}
	}()
	waitForShutdown()
}

// initEngine 启动撮合 临时
func initEngine(cache *redis.BooksCache) {
	logger.Info("init engine")
	// todo 交易对处理
	coinParis := []string{"btc_usdt", "btc_eth", "btc_bnb"}
	for _, coinPair := range coinParis {
		engine := internal.NewEngine(cache, coinPair)
		go engine.Start(coinPair)
		engines[coinPair] = engine
		logger.Infof("init engine for %s success", coinPair)
	}
}

// 输出当前的交易列表
func printOrderBook(engine *internal.Engine) {
	ticker := time.NewTicker(20 * time.Second)
	for {
		<-ticker.C
		engine.GetOrderBookList()
	}
}

func waitForShutdown() {
	sign := make(chan os.Signal, 1)
	signal.Notify(sign, ShutdownSignals...)

	select {
	case <-sign:
		time.AfterFunc(time.Minute*5, func() {
			logger.Warn("WaitForShutdown time out")
			os.Exit(-1)
		})

		for _, engine := range engines {
			engine.Shutdown()
		}
		logger.Infof("WaitForShutdown success")
		os.Exit(0)
	}
}
