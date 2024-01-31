package config

import (
	"github.com/go-micro/plugins/v4/logger/zap"
	"go-micro.dev/v4/logger"
	z "go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitLogger() {
	config := z.NewProductionConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.OutputPaths = []string{Conf.Log.Info}
	l, _ := config.Build()

	withLogger := zap.WithLogger(l)

	newLogger, err := zap.NewLogger(withLogger)
	if err != nil {
		panic(err)
	}

	logger.DefaultLogger = newLogger
}
