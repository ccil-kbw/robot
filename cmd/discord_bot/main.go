// Package main
package main

import (
	"flag"
	"fmt"

	"github.com/ccil-kbw/robot/discord"
)

var (
	GuildID        = flag.String("guild", "", "Test guild ID. If not passed - bot registers commands globally")
	BotToken       = flag.String("token", "", "Bot access token")
	RemoveCommands = flag.Bool("rmcmd", true, "Remove all commands after shutdowning or not")
)

func init() { flag.Parse() }

func main() {
	msgs := make(chan string)
	go discord.Run(GuildID, BotToken, RemoveCommands, msgs)

	for msg := range msgs {
		fmt.Println(msg)
	}
}
