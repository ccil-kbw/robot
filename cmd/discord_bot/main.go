package main

import (
	"log"

	"github.com/caarlos0/env/v11"
	config "github.com/ccil-kbw/robot/internal/config"
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
	bot := discord.NewDiscordBot(logger.Logger, config.Cfg.DiscordServerID, config.Cfg.DiscordBotToken, config.Cfg.DiscordBotAsPublic)
	bot.StartBot()
}

func loadConfig() {
	if err := env.Parse(&config.Cfg); err != nil {
		log.Fatalf("%+v\n", err)
	}
	log.Println("configuration ready")
}
