package main

import (
	"fmt"
	iqamav2 "github.com/ccil-kbw/robot/pkg/iqama/v2"
	"github.com/ccil-kbw/robot/pkg/masjid_info"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: iqamacli <@masjid>, e.g: iqamacli @ccil_kbw")
		os.Exit(1)
	}
	masjidFolderRoot := os.Args[1]
	client := iqamav2.NewIqamaCSV(masjidFolderRoot)

	fmt.Println(client.GetShellPrettified())

	masjidInfo := masjid_info.GetMasjidInfoFromFile(masjidFolderRoot)
	fmt.Println(masjidInfo)
}
