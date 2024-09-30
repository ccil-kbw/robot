package main

import (
	"fmt"

	environment "github.com/ccil-kbw/robot/internal/environment"
	"github.com/ccil-kbw/robot/pkg/discord"
	rec2 "github.com/ccil-kbw/robot/pkg/rec"

	"os"
	"os/signal"
	"strings"
	"time"

	"go.uber.org/zap"
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
	environment.LoadEnvironmentVariables()
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

		go func() {
			for {
				data := rec2.NewRecordConfigDataS()
				obsClient, err := rec2.StartRecServer(host, password, data)
				if err != nil {
					fmt.Printf("could not reach or authenticate to OBS, retrying in 1 minute...\n")
					time.Sleep(1 * time.Minute)
					continue
				}

				// Check recording status every minute
				for {
					data = rec2.NewRecordConfigDataS()
					shouldRecord := rec2.SupposedToBeRecording(data)
					isRecording, err := obsClient.IsRecording()
					if err != nil {
						fmt.Printf("couldn't check if OBS is recording: %v\n", err)
						continue
					}

					if shouldRecord && !isRecording {
						fmt.Println("should be recording")
						err := obsClient.StartRecording()
						if err != nil {
							fmt.Printf("couldn't start recording: %v\n", err)
						}
					} else if !shouldRecord && isRecording {
						fmt.Println("should not be recording")
						err := obsClient.StopRecording()
						if err != nil {
							fmt.Printf("couldn't stop recording: %v\n", err)
						}
					}
					time.Sleep(1 * time.Minute)

				}
			}
		}()
	}

out:
	for {
		select {
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

	for {
		err := discordBot.StartBot()
		if err != nil {
			logger.Error("Failed to start Discord bot, retrying in 10 minutes", zap.Error(err))
			time.Sleep(10 * time.Minute)
		} else {
			break
		}
	}
}
