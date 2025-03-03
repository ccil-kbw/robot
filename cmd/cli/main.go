package main

import (
	"fmt"
	iqamav2 "github.com/ccil-kbw/robot/iqama/v2"
)

func main() {
	client := iqamav2.NewIqamaCSV("iqama_2025.csv")

	fmt.Println(client.GetShellPrettified())
}
