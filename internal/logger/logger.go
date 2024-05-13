package logger

import (
	"log"

	"go.uber.org/zap"
)

var (
	Logger *zap.Logger
	Cfg    = config{}
)

type config struct {
	Environment        string `env:"ENVIRONMENT" envDefault:"dev"`
	DiscordServerID    string `env:"DISCORD_SERVER_ID"`
	DiscordBotToken    string `env:"DISCORD_BOT_TOKEN"`
	DiscordBotAsPublic bool   `env:"DISCORD_BOT_AS_PUBLIC" envDefault:"true"`
}

func InitializeLogger() {
	var err error
	if Cfg.Environment == "dev" {
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
