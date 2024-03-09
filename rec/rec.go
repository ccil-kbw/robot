package rec

import (
	"fmt"

	"github.com/andreykaipov/goobs"
)

type Recorder struct {
	client *goobs.Client
}

func (o *Recorder) GetClient() *goobs.Client {
	return o.client
}

func New(host, password string) (*Recorder, error) {

	if host == "" {
		host = "localhost:4455"
	}

	client, err := goobs.New(host, goobs.WithPassword(password))
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

func (o *Recorder) DispatchOperation(msg string) error {
	fmt.Printf("[todo] dispatching operation: %s\n", msg)
	return nil
}

func (o *Recorder) Disconnect() error {
	return o.client.Disconnect()
}

// StartRecording records if didn't start recording, else continue recording for the specified duration
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

func (o *Recorder) RecordTime() float64 {
	r, err := o.client.Record.GetRecordStatus()
	if err != nil {
		fmt.Println(err)
	}
	return r.OutputDuration
}
