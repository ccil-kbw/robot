package rec

import (
	"fmt"
)

func StartRecServer(host, password string, data *RecordConfigDataS) (*Recorder, error) {
	client, err := New(host, password)
	if err != nil {
		fmt.Println("could not initiate client")
		return nil, err
	}

	fmt.Println("Starting OBS recording control routine")
	return client, nil
}
