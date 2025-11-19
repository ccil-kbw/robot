package v1

import (
	"encoding/json"
	"fmt"
	"os"

	v2 "github.com/ccil-kbw/robot/iqama/v2"
)

// Get returns today's prayer times in v1 format by converting from v2 CSV data
func Get() (*Resp, error) {
	// Get CSV file path from environment or use default
	csvPath := os.Getenv("MDROID_IQAMA_CSV_PATH")
	if csvPath == "" {
		csvPath = "iqama_2025.csv"
	}

	// Use v2 to read CSV data
	iqamaClient := v2.NewIqamaCSV(csvPath)
	todayTimes, err := iqamaClient.GetTodayTimes()
	if err != nil {
		return nil, fmt.Errorf("failed to get today's times: %w", err)
	}

	// Convert v2 format to v1 format
	resp := &Resp{
		Dates: Dates{
			Hijri:  "", // v2 doesn't have Hijri dates
			Miladi: todayTimes.Date.Format("01/02/2006"),
		},
		Fajr: Prayer{
			Adhan: todayTimes.Fajr.Adhan.Format("15:04"),
			Iqama: todayTimes.Fajr.Iqama.Format("15:04"),
		},
		Sunrise: "", // v2 doesn't track sunrise
		Dhuhr: Prayer{
			Adhan: todayTimes.Dhuhr.Adhan.Format("15:04"),
			Iqama: todayTimes.Dhuhr.Iqama.Format("15:04"),
		},
		Asr: Prayer{
			Adhan: todayTimes.Asr.Adhan.Format("15:04"),
			Iqama: todayTimes.Asr.Iqama.Format("15:04"),
		},
		Maghrib: Prayer{
			Adhan: todayTimes.Maghrib.Adhan.Format("15:04"),
			Iqama: todayTimes.Maghrib.Iqama.Format("15:04"),
		},
		Isha: Prayer{
			Adhan: todayTimes.Isha.Adhan.Format("15:04"),
			Iqama: todayTimes.Isha.Iqama.Format("15:04"),
		},
		Jumua: Jumua{
			Fr: "", // v2 doesn't have Jumua times
			Ar: "",
		},
	}

	return resp, nil
}

// GetRAW returns today's prayer times as JSON bytes
func GetRAW() []byte {
	resp, err := Get()
	if err != nil {
		// Return error as JSON instead of panicking
		errResp := map[string]string{"error": err.Error()}
		data, _ := json.Marshal(errResp)
		return data
	}

	data, err := resp.Marshal()
	if err != nil {
		errResp := map[string]string{"error": fmt.Sprintf("failed to marshal response: %v", err)}
		data, _ := json.Marshal(errResp)
		return data
	}

	return data
}
