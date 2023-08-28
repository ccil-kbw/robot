package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/ccil-kbw/robot/discord"
	v1 "github.com/ccil-kbw/robot/iqama/v1"
)

var (
	// First iteration notes:
	// Hardcode to enable the feature on the code side
	// Double check on the env side if we need the feature at run time (e.g openbroadcaster is configured, proxy is required, etc)
	config = Config{
		Features: Features{
			Proxy:      true,
			DiscordBot: true,
		},
	}
)

type Features struct {
	Proxy      bool
	DiscordBot bool
	Record     bool
}

type Config struct {
	Features Features
}

func main() {
	if config.Features.Proxy {
		go proxy()
	}

	if config.Features.DiscordBot {
		go bot()
	}

	// handle erroring, for now just block
	<-make(chan struct{})
}

// proxy, move to apis, maybe pkg/apis/proxyserver/proxyserver.go
func proxy() {
	http.HandleFunc("/today", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("request: %s, %s\n", r.Method, r.URL)
		io.WriteString(w, string(v1.GetRAW()))
	})

	fmt.Println("Running iqama-proxy Go server on port :3333")
	_ = http.ListenAndServe(":3333", nil)
}

func bot() {
	guildID := os.Getenv("MDROID_BOT_GUILD_ID")
	botToken := os.Getenv("MDROID_BOT_TOKEN")
	removeCommands := true

	discord.Run(&guildID, &botToken, &removeCommands)
}
