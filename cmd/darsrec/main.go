package main

import (
	"fmt"
	"github.com/ccil-kbw/robot/iqama"
	"os"
	"time"

	"github.com/ccil-kbw/robot/rec"
)

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

	data := iqama.StartRecordingScheduleServer()

	for {

		isRecording, err := client.IsRecording()
		if err != nil {
			panic(err)
		}

		shouldRecord := rec.SupposedToBeRecording(data.Confs())

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
