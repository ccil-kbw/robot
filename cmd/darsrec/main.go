package main

import (
	"fmt"
	"os"

	environment "github.com/ccil-kbw/robot/internal/environment"
	rec2 "github.com/ccil-kbw/robot/pkg/rec"
)

func init() {
	environment.LoadEnvironmentVariables()
}

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
	data := rec2.NewRecordConfigDataS()

	rec2.StartRecServer(host, password, data)

}
