package main

import (
	"log"

	"github.com/caarlos0/env/v11"
	environment "github.com/ccil-kbw/robot/internal/environment"
	logger "github.com/ccil-kbw/robot/internal/logger"
	"github.com/ccil-kbw/robot/pkg/discord"
)

func init() {
	environment.LoadEnvironmentVariables()
	logger.InitializeLogger()
	loadConfig()
}

func main() {
	bot := discord.NewDiscordBot(logger.Logger, logger.Cfg.DiscordServerID, logger.Cfg.DiscordBotToken, logger.Cfg.DiscordBotAsPublic)
	bot.StartBot()
}

func loadConfig() {
	if err := env.Parse(&logger.Cfg); err != nil {
		log.Fatalf("%+v\n", err)
	}
	log.Println("configuration ready")
}
