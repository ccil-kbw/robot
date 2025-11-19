// Package main
package main

import (
	"flag"
	"log"

	"github.com/ccil-kbw/robot/discord"
)

var (
	GuildID        = flag.String("guild", "", "Test guild ID. If not passed - bot registers commands globally")
	BotToken       = flag.String("token", "", "Bot access token")
	RemoveCommands = flag.Bool("rmcmd", true, "Remove all commands after shutdowning or not")
)

func init() { flag.Parse() }

func main() {
	if err := discord.Run(GuildID, BotToken, RemoveCommands, nil, make(chan string)); err != nil {
		log.Fatalf("Discord bot failed: %v", err)
	}
}
