package main

import (
	"log"

	"github.com/caarlos0/env/v11"
	environment "github.com/ccil-kbw/robot/internal/environment"
	"github.com/ccil-kbw/robot/pkg/discord"
	"go.uber.org/zap"
)

// TODO: Should Types have a folder of their own?
type config struct {
	Environment        string `env:"ENVIRONMENT" envDefault:"dev"`
	DiscordServerID    string `env:"DISCORD_SERVER_ID"`
	DiscordBotToken    string `env:"DISCORD_BOT_TOKEN"`
	DiscordBotAsPublic bool   `env:"DISCORD_BOT_AS_PUBLIC" envDefault:"true"`
}

var (
	cfg    = config{}
	logger *zap.Logger
)

func init() {
	environment.LoadEnvironmentVariables()
	loadConfig()
	initializeLogger()
}

func main() {
	bot := discord.NewDiscordBot(logger, cfg.DiscordServerID, cfg.DiscordBotToken, cfg.DiscordBotAsPublic)
	bot.StartBot()
}

func loadConfig() {
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("%+v\n", err)
	}
	log.Println("configuration ready")
}

func initializeLogger() {
	var err error
	if cfg.Environment == "dev" {
		logger, err = zap.NewDevelopment()
		if err != nil {
			initExampleLoggerWhenFail()
		}
		logger.Info("Using Development Logger")
	} else {
		logger, err = zap.NewProduction()
		if err != nil {
			initExampleLoggerWhenFail()
		}
		logger.Info("Using Production Logger")
	}
}

func initExampleLoggerWhenFail() {
	log.Println("can't initialize zap logger, using basic logger")
	logger = zap.NewExample()
}
