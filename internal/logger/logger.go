package logger

import (
	"log"

	config "github.com/ccil-kbw/robot/internal/config"
	"go.uber.org/zap"
)

var (
	Logger *zap.Logger
)

func InitializeLogger() {
	var err error
	if config.Cfg.Environment == "dev" {
		Logger, err = zap.NewDevelopment()
		if err != nil {
			initExampleLoggerWhenFail()
		}
		Logger.Info("Using Development Logger")
	} else {
		Logger, err = zap.NewProduction()
		if err != nil {
			initExampleLoggerWhenFail()
		}
		Logger.Info("Using Production Logger")
	}
}

func initExampleLoggerWhenFail() {
	log.Println("can't initialize zap logger, using basic logger")
	Logger = zap.NewExample()
}
