package masjid_info

import (
	"fmt"
	"log"
	"os"

	yaml "gopkg.in/yaml.v2"
)

type Coordinates struct {
	Latitude  float64 `yaml:"latitude"`
	Longitude float64 `yaml:"longitude"`
}

type MasjidInfo struct {
	Name          string      `yaml:"name"`
	Address       string      `yaml:"address"`
	Website       string      `yaml:"website"`
	Coordinates   Coordinates `yaml:"coordinates"`
	ContributedBy string      `yaml:"contributed_by"`
}

func GetMasjidInfoFromFile(folderRoot string) MasjidInfo {
	var info MasjidInfo

	// Read the YAML file
	file, err := os.ReadFile(fmt.Sprintf("assets/masjids_data/%s/info.yaml", folderRoot))
	if err != nil {
		log.Fatalf("Failed to read file: %s", err)
	}

	// Unmarshal the YAML file into the MasjidInfo struct
	err = yaml.Unmarshal(file, &info)
	if err != nil {
		log.Fatalf("Failed to unmarshal YAML: %s", err)
	}

	return info
}
