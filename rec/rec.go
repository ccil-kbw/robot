package rec

import (
	"fmt"

	"github.com/andreykaipov/goobs"
)

type Recorder struct {
	client *goobs.Client
}

func New(password string) (*Recorder, error) {

	client, err := goobs.New("localhost:4455", goobs.WithPassword(password))
	if err != nil {
		return nil, err
	}

	version, err := client.General.GetVersion()
	if err != nil {
		return nil, err
	}

	fmt.Printf("OBS Studio version: %s\n", version.ObsVersion)
	fmt.Printf("Websocket server version: %s\n", version.ObsWebSocketVersion)

	fmt.Println("Connected to OBS Studio. Starting Routines")

	return &Recorder{
		client: client,
	}, nil
}

func (o *Recorder) Disconnect() error {
	return o.client.Disconnect()
}

// Rec records if didn't start recording, else continue recording for the specified duration
func (o *Recorder) StartRecording() error {

	// Get Record Status
	recordStatus, err := o.client.Record.GetRecordStatus()
	if err != nil {
		return err
	}

	if recordStatus.OutputActive {
		fmt.Printf("Record in Progress. Duration: %f seconds\n", recordStatus.OutputDuration/1000)
	}

	if !recordStatus.OutputActive {
		fmt.Println("Starting record")
		o.client.Record.StartRecord()
	}

	return nil
}

func (o *Recorder) StopRecording() error {
	// Get Record Status
	recordStatus, err := o.client.Record.GetRecordStatus()
	if err != nil {
		return err
	}

	if recordStatus.OutputActive {
		fmt.Printf("Recorded %f seconds clip. Saving...\n", recordStatus.OutputDuration/1000)
		o.client.Record.StopRecord()
	}

	if !recordStatus.OutputActive {
		fmt.Println("Not recording, nothing to stop.")
	}

	return nil
}

func (o *Recorder) IsRecording() (bool, error) {
	resp, err := o.client.Record.GetRecordStatus()
	if err != nil {
		return false, err
	}

	return resp.OutputActive, nil
}
