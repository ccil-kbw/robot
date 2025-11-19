package main

import (
	"fmt"
	"log"

	iqamav2 "github.com/ccil-kbw/robot/iqama/v2"
)

func main() {
	client, err := iqamav2.NewIqamaCSV(iqamav2.GetDefaultIqamaCSVPath())
	if err != nil {
		log.Fatalf("Failed to initialize iqama client: %v", err)
	}

	fmt.Println(client.GetShellPrettified())
}
