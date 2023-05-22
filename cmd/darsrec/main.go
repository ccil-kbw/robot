package main

import (
	"fmt"
	"os"
	"time"

	"github.com/ccil-kbw/robot/rec"
)

func main() {

	password := os.Getenv("OBS_WEBSOCKET_PASSWORD")
	if password == "" {
		fmt.Println("OBS_WEBSOCKET_PASSWORD is unset, cannot proceed")
		os.Exit(1)
	}

	client, err := rec.New(password)
	if err != nil {
		panic(err)
	}

	defer client.Disconnect()

	confs := rec.GetIqamaRecordingConfigs()

	for {
		isRecording, err := client.IsRecording()
		if err != nil {
			panic(err)
		}

		shouldRecord := rec.SupposedToBeRecording(confs)

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
