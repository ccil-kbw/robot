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
			in := 15 * time.Minute
			notifyFunc(notifyChan, prayersData, in)
			notifyFunc(notifyChan, prayersData, 0)
			time.Sleep(10 * time.Second)

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
		_, _ = io.WriteString(w, string(v1.GetRAW()))
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

func notifyFunc(notifyChan chan string, prayersData *iqama.PrayersData, in time.Duration) {
	now := time.Now().Add(in)

	if prayersData.Confs().Fajr.Iqama == now.Format(fmt.Sprintf("%s:%s", stdHour, stdMinutes)) {
		notifyPrayer("Fajr", prayersData.Confs().Fajr.Iqama, in, notifyChan)
	}

	if prayersData.Confs().Dhuhr.Iqama == now.Format(fmt.Sprintf("%s:%s", stdHour, stdMinutes)) {
		notifyPrayer("Dhuhr", prayersData.Confs().Dhuhr.Iqama, in, notifyChan)

	}

	if prayersData.Confs().Asr.Iqama == now.Format(fmt.Sprintf("%s:%s", stdHour, stdMinutes)) {
		notifyPrayer("Asr", prayersData.Confs().Asr.Iqama, in, notifyChan)
	}

	if prayersData.Confs().Maghrib.Iqama == now.Format(fmt.Sprintf("%s:%s", stdHour, stdMinutes)) {
		notifyPrayer("Maghrib", prayersData.Confs().Maghrib.Iqama, in, notifyChan)
	}

	if prayersData.Confs().Isha.Iqama == now.Format(fmt.Sprintf("%s:%s", stdHour, stdMinutes)) {
		notifyPrayer("Isha", prayersData.Confs().Isha.Iqama, in, notifyChan)
	}

}

func notifyPrayer(prayerName, prayerTime string, in time.Duration, notifyChan chan string) {
	var msg string
	{
		if in == 0 {
			msg = fmt.Sprintf("%s's Iqama Time now, it's %s!", prayerName, prayerTime)
		} else {
			msg = fmt.Sprintf("%s's Iqama in %v, at %s", prayerName, in, prayerTime)
		}
	}
	notifyChan <- msg
}
