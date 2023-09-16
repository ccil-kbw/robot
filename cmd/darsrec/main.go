package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/ccil-kbw/robot/rec"
)

// AsyncConf (sync.Mutex) is used to handle the RecordConfig asynchronously as it's used in the main and some sub routines
type AsyncConf struct {
	mu    sync.Mutex
	confs []*rec.RecordConfig
}

// Refresh updates the scheduling configurations
func (ac *AsyncConf) Refresh() {
	ac.mu.Lock()
	ac.confs = rec.GetIqamaRecordingConfigs()
	ac.mu.Unlock()
}

func (ac *AsyncConf) Confs() []*rec.RecordConfig {
	ac.mu.Lock()
	defer ac.mu.Unlock()
	return ac.confs
}

func main() {

	var host, password string
	{
		host = os.Getenv("OBS_WEBSOCKET_HOST")
		if host == "" {
			host = "localhost:4455"
			fmt.Printf("OBS_WEBSOCKET_HOST not set, using default %s\f", host)
		}

		password = os.Getenv("OBS_WEBSOCKET_PASSWORD")
		if password == "" {
			fmt.Println("OBS_WEBSOCKET_PASSWORD is unset, cannot proceed")
			os.Exit(1)
		}
	}

	client, err := rec.New(host, password)
	if err != nil {
		panic(err)
	}

	defer client.Disconnect()

	ac := AsyncConf{}
	ac.Refresh()

	ticker := time.NewTicker(1 * time.Hour)
	go func() {
		for range ticker.C {
			fmt.Printf("%s: Updating Iqama time...\n", time.Now().Format(time.Kitchen))
			ac.Refresh()
		}
	}()

	for {

		isRecording, err := client.IsRecording()
		if err != nil {
			panic(err)
		}

		shouldRecord := rec.SupposedToBeRecording(ac.Confs())

		// Start recording if supposed to be recording but currently not recording
		if shouldRecord && !isRecording {
			err := client.StartRecording()
			if err != nil {
				panic(err)
			}
		}

		// Stop recording if not supposed to be recording but currently recording
		if !shouldRecord && isRecording {
			err := client.StopRecording()
			if err != nil {
				panic(err)
			}
		}

		time.Sleep(1 * time.Minute)
	}
}
