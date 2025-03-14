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

var (
	stdHour    = "15"
	stdMinutes = "4"
)

func main() {
	msgs := make(chan string)
	stop := make(chan os.Signal, 1)
	notifyChan := make(chan string)

	signal.Notify(stop, os.Interrupt)
	/*
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
				time.Sleep(55 * time.Second)
			}

		}()
	*/
	var obs *rec.Recorder

	if config.Features.DiscordBot {
		go bot(obs, notifyChan)
	}

	if config.Features.Record {
		host := os.Getenv("MDROID_OBS_WEBSOCKET_HOST")
		password := os.Getenv("MDROID_OBS_WEBSOCKET_PASSWORD")

		obsClient := startServerWithRetry(host, password)

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
				obsClient = startServerWithRetry(host, password)
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

func startServerWithRetry(host string, password string) *rec.Recorder {
	for {
		obs, err := rec.StartRecServer(host, password)
		if err != nil {
			fmt.Printf("could not reach or authenticate to OBS, retrying in 1 minutes...\n")
			time.Sleep(1 * time.Minute)
		} else {
			return obs
		}
	}
}
