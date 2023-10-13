package main

import (
	"fmt"
	"os"
	"time"

	"github.com/ccil-kbw/robot/rec"
)

// TODO: Swap all this as a darsRec.StartServer()
func main() {

	var host, password string
	{
		host = os.Getenv("OBS_WEBSOCKET_HOST")
		if host == "" {
			host = "localhost:4455"
			fmt.Printf("OBS_WEBSOCKET_HOST not set, using default %s\n", host)
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

	defer func(client *rec.Recorder) {
		err := client.Disconnect()
		if err != nil {
			fmt.Println("couldn't disconnect client")
		}
	}(client)

	for {

		isRecording, err := client.IsRecording()
		if err != nil {
			fmt.Println("couldn't check if OBS is recording")
		}

		shouldRecord := rec.SupposedToBeRecording()

		// Start recording if supposed to be recording but currently not recording
		if shouldRecord && !isRecording {
			err := client.StartRecording()
			if err != nil {
				fmt.Println("couldn't start recording")
			}
		}

		var recordTimeLimit float64
		{
			recordTimeLimit = 2 * 60 * 60 * 1000
		}

		// Stop recording if not supposed to be recording but currently recording
		if !shouldRecord && isRecording && (client.RecordTime() > recordTimeLimit) {
			err := client.StopRecording()
			if err != nil {
				fmt.Println("couldn't stop recording")
			}
		}

		time.Sleep(1 * time.Minute)
	}
}
