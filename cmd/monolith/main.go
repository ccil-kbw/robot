package main

import (
	"fmt"
	"github.com/ccil-kbw/robot/iqama"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/ccil-kbw/robot/discord"
	v1 "github.com/ccil-kbw/robot/iqama/v1"
	"github.com/ccil-kbw/robot/rec"
)

var (
	// First iteration notes:
	// Hardcode to enable the feature on the code side
	// Double check on the env side if we need the feature at run time (e.g openbroadcaster is configured, proxy is required, etc)
	config = Config{
		Features: Features{
			Proxy:      true,
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

func main() {
	var err error
	msgs := make(chan string)
	stop := make(chan os.Signal, 1)
	notifyChan := make(chan string)

	signal.Notify(stop, os.Interrupt)

	var prayersData *iqama.PrayersData
	{
		prayersData = iqama.StartIqamaServer()
	}

	go func() {
		go iqama.StartRecordingScheduleServer()

		for {
			now := time.Now()
			if prayersData.Confs().Fajr.Iqama == now.Format(fmt.Sprintf("%s:%s", stdHour, stdMinutes)) {
				fmt.Println("Fajr: Iqama starts Now")
				notifyChan <- "Fajr: Iqama starts Now"
			}

			if prayersData.Confs().Dhuhr.Iqama == now.Format(fmt.Sprintf("%s:%s", stdHour, stdMinutes)) {
				fmt.Println("Dhuhur: Iqama Starts Now")
				notifyChan <- "Dhuhur: Iqama starts Now"

			}

			if prayersData.Confs().Asr.Iqama == now.Format(fmt.Sprintf("%s:%s", stdHour, stdMinutes)) {
				fmt.Println("Asr: Iqama Starts Now")
				notifyChan <- "Asr: Iqama Starts Now"
			}

			if prayersData.Confs().Maghrib.Iqama == now.Format(fmt.Sprintf("%s:%s", stdHour, stdMinutes)) {
				fmt.Println("Maghrib: Iqama Starts Now")
				notifyChan <- "Maghrib: Iqama Starts Now"
			}

			if prayersData.Confs().Isha.Iqama == now.Format(fmt.Sprintf("%s:%s", stdHour, stdMinutes)) {
				fmt.Println("Isha: Iqama Starts Now")
				notifyChan <- "Isha: Iqama Starts Now"
			}

			time.Sleep(1 * time.Minute)
		}

	}()

	if config.Features.Proxy {
		go proxy()
	}

	var obs *rec.Recorder

	if config.Features.Record {
		host := os.Getenv("MDROID_OBS_WEBSOCKET_HOST")
		password := os.Getenv("MDROID_OBS_WEBSOCKET_PASSWORD")
		obs, err = rec.New(host, password)
		if err != nil {
			fmt.Printf("could not reach or authenticate to OBS")
		}
	}

	if config.Features.DiscordBot {
		go bot(obs, notifyChan)
	}

out:
	for {
		select {
		// discord msgs dispatcher
		case msg := <-msgs:
			fmt.Printf("%v, operation received from discord: %s\n", time.Now(), msg)
			if strings.HasPrefix(msg, "obs-") {
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

// proxy, move to apis, maybe pkg/apis/proxyserver/proxyserver.go
func proxy() {
	http.HandleFunc("/today", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("request: %s, %s\n", r.Method, r.URL)
		io.WriteString(w, string(v1.GetRAW()))
	})

	fmt.Println("Running iqama-proxy Go server on port :3333")
	_ = http.ListenAndServe(":3333", nil)
}

func bot(obs *rec.Recorder, notifyChan chan string) {
	guildID := os.Getenv("MDROID_BOT_GUILD_ID")
	botToken := os.Getenv("MDROID_BOT_TOKEN")
	removeCommands := true
	discord.Run(&guildID, &botToken, &removeCommands, obs, notifyChan)
}
