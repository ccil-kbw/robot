package main

import (
	"fmt"
	iqamav2 "github.com/ccil-kbw/robot/iqama/v2"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: iqamacli <@masjid>, e.g: iqamacli @ccil_kbw")
		os.Exit(1)
	}
	client := iqamav2.NewIqamaCSV(os.Args[1])

	fmt.Println(client.GetShellPrettified())
}
