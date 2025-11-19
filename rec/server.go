package rec

import (
	"fmt"
	"time"
)

// StartRecServer initializes an OBS WebSocket client and starts an automatic recording control routine.
// It connects to OBS at the specified host using the provided password.
// The routine runs in a goroutine and checks every minute whether recording should be active
// based on the current prayer schedule. It automatically starts or stops recording as needed.
// Returns the Recorder client for manual control, or an error if connection fails.
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
