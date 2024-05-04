package main

import (
	"fmt"
	"log"

	"github.com/ccil-kbw/robot/pkg/discord"
	rec2 "github.com/ccil-kbw/robot/pkg/rec"

	"os"
	"os/signal"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/joho/godotenv"
)

var (
	// First iteration notes:
	// Hardcode to enable the feature on the code side
	// Double check on the env side if we need the feature at run time (e.g openbroadcaster is configured, proxy is required, etc)
	config = Config{
		Features: Features{
			DiscordBot: true,
			Record:     true,
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

var (
	stdHour    = "15"
	stdMinutes = "4"
)

func init() {
	loadEnvs()
}

func main() {
	msgs := make(chan string)
	stop := make(chan os.Signal, 1)

	signal.Notify(stop, os.Interrupt)

	var obs *rec2.Recorder

	if config.Features.DiscordBot {
		go bot()
	}

	if config.Features.Record {
		host := os.Getenv("MDROID_OBS_WEBSOCKET_HOST")
		password := os.Getenv("MDROID_OBS_WEBSOCKET_PASSWORD")
		data := rec2.NewRecordConfigDataS()

		obsClient := startServerWithRetry(host, password, data)

		// Calculate the duration until midnight
		now := time.Now()
		night := time.Date(now.Year(), now.Month(), now.Day()+1, 2, 0, 0, 0, now.Location())
		duration := night.Sub(now)

		// Create a timer that waits until midnight
		timer := time.NewTimer(duration)
		<-timer.C // This blocks until the timer fires

		// Now that it's midnight, start a ticker that ticks every 24 hours
		ticker := time.NewTicker(24 * time.Hour)

		// Call rec.StartRecServer every time the ticker ticks
		go func() {
			for range ticker.C {
				err := obsClient.Disconnect()
				if err != nil {
					return
				}
				obsClient = startServerWithRetry(host, password, data)
			}
		}()

	}

out:
	for {
		select {
		// discord msgs dispatcher
		case msg := <-msgs:
			fmt.Printf("%v, operation received from discord: %s\n", time.Now(), msg)
			if strings.HasPrefix(msg, "rec-") {
				if config.Features.Record {
					fmt.Println("feature enabled")
					err := obs.DispatchOperation(msg)
					if err != nil {
						fmt.Println("failed dispatching operation")
					}
				}
			}
		case <-stop:
			// simplest way to wait for the nested go routines to clean up
			// takes < 2 ms but better be safe
			time.Sleep(10 * time.Second)
			break out
		}
	}

}

func bot() {
	guildID := os.Getenv("MDROID_BOT_GUILD_ID")
	botToken := os.Getenv("MDROID_BOT_TOKEN")

	logger, _ := zap.NewProduction()
	discordBot := discord.NewDiscordBot(logger, guildID, botToken, true)
	discordBot.StartBot()
}

func startServerWithRetry(host string, password string, data *rec2.RecordConfigDataS) *rec2.Recorder {
	for {
		obs, err := rec2.StartRecServer(host, password, data)
		if err != nil {
			fmt.Printf("could not reach or authenticate to OBS, retrying in 1 minutes...\n")
			time.Sleep(1 * time.Minute)
		} else {
			return obs
		}
	}
}

func loadEnvs() {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env not present, using process envs and defaults")
	} else {
		log.Println(".env loaded")
	}
}
