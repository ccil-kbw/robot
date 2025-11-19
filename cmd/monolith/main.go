package main

import (
	"fmt"
	"github.com/ccil-kbw/robot/iqama"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/ccil-kbw/robot/discord"
	"github.com/ccil-kbw/robot/rec"
)

var (
	// First iteration notes:
	// Hardcode to enable the feature on the code side
	// Double check on the env side if we need the feature at run time (e.g openbroadcaster is configured, proxy is required, etc)
	config = Config{
		Features: Features{
			DiscordBot: false,
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


func main() {
	msgs := make(chan string)
	stop := make(chan os.Signal, 1)
	notifyChan := make(chan string)

	signal.Notify(stop, os.Interrupt)

	var obs *rec.Recorder

	if config.Features.DiscordBot {
		go bot(obs, notifyChan)
	}

	if config.Features.Record {
		host := os.Getenv("MDROID_OBS_WEBSOCKET_HOST")
		password := os.Getenv("MDROID_OBS_WEBSOCKET_PASSWORD")

		obsClient := startServerWithRetry(host, password)
		if obsClient == nil {
			fmt.Println("Warning: Failed to initialize OBS client. Recording features will be disabled.")
		} else {
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
					if obsClient != nil {
						if err := obsClient.Disconnect(); err != nil {
							fmt.Printf("Error disconnecting from OBS: %v\n", err)
						}
					}
					newClient := startServerWithRetry(host, password)
					if newClient != nil {
						obsClient = newClient
					} else {
						fmt.Println("Failed to reconnect to OBS")
					}
				}
			}()
		}
		obs = obsClient
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

func bot(obs *rec.Recorder, notifyChan chan string) {
	guildID := os.Getenv("MDROID_BOT_GUILD_ID")
	botToken := os.Getenv("MDROID_BOT_TOKEN")
	removeCommands := true
	if err := discord.Run(&guildID, &botToken, &removeCommands, obs, notifyChan); err != nil {
		fmt.Printf("Discord bot error: %v\n", err)
	}
}

func startServerWithRetry(host string, password string) *rec.Recorder {
	maxRetries := 10
	initialDelay := 10 * time.Second
	maxDelay := 5 * time.Minute

	for attempt := 1; attempt <= maxRetries; attempt++ {
		obs, err := rec.StartRecServer(host, password)
		if err == nil {
			return obs
		}

		if attempt == maxRetries {
			fmt.Printf("Failed to connect to OBS after %d attempts. Giving up.\n", maxRetries)
			return nil
		}

		// Exponential backoff: 10s, 20s, 40s, 80s, 160s (capped at 5min)
		delay := initialDelay * time.Duration(1<<uint(attempt-1))
		if delay > maxDelay {
			delay = maxDelay
		}

		fmt.Printf("Could not reach or authenticate to OBS (attempt %d/%d): %v\n", attempt, maxRetries, err)
		fmt.Printf("Retrying in %v...\n", delay)
		time.Sleep(delay)
	}

	return nil
}
