package rec

import (
	"fmt"
	"time"
)

func StartRecServer(host, password string, data *RecordConfigDataS) (*Recorder, error) {
	client, err := New(host, password)
	if err != nil {
		fmt.Println("could not initiate client")
		return nil, nil
	}

	fmt.Println("Starting OBS recording control routine")
	go func() {
		for {

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
	}()

	return client, nil
}
