package rec

import (
	"fmt"
	"time"
)

func StartRecServer(host, password string) (*Recorder, error) {
	client, err := New(host, password)
	if err != nil {
		fmt.Println("could not initiate client")
		return nil, err
	}

	fmt.Println("Starting OBS recording control routine")
	go func() {
		for {
			data := NewRecordConfigDataS()

			isRecording, err := client.IsRecording()
			if err != nil {
				fmt.Println("couldn't check if OBS is recording")
			}

			shouldRecord := SupposedToBeRecording(data)

			// Start recording if supposed to be recording but currently not recording
			if shouldRecord && !isRecording {
				err := client.StartRecording()
				if err != nil {
					fmt.Println("couldn't start recording")
				}
			}

			// Stop recording if not supposed to be recording but currently recording
			if !shouldRecord && isRecording {
				err := client.StopRecording()
				if err != nil {
					fmt.Println("couldn't stop recording")
				}
			}

			time.Sleep(1 * time.Minute)
		}
	}()

	return client, nil
}
