package main

import (
	"github.com/caarlos0/env/v11"
	"github.com/ccil-kbw/robot/pkg/discord"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"log"
)

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
	loadEnvs()
	loadConfig()
	initializeLogger()
}

func main() {
	bot := discord.NewDiscordBot(logger, cfg.DiscordServerID, cfg.DiscordBotToken, cfg.DiscordBotAsPublic)
	bot.StartBot()
}

func loadEnvs() {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env not present, using process envs and defaults")
	} else {
		log.Println(".env loaded")
	}
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
