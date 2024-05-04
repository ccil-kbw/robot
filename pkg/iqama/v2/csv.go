package v2

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"time"
)

type IqamaCSV struct {
	filePath   string
	iqamaTimes map[string]IqamaDailyTimes
}

func NewIqamaCSV(masjid string) Iqama {
	filePath := fmt.Sprintf("assets/masjids_data/%s/iqama.csv", masjid)
	i := &IqamaCSV{filePath: filePath}
	if err := i.readCSV(); err != nil {
		log.Fatalf("Unable to read CSV file %s", filePath)
	}
	return i
}

func (i *IqamaCSV) GetTodayTimes() (*IqamaDailyTimes, error) {
	// Get today's date
	today := time.Now()

	times, err := i.iqamaForDate(today)
	handleErr(err)
	return &times, nil
}

func (i *IqamaCSV) GetTomorrowTimes() (*IqamaDailyTimes, error) {
	// Get tomorrow's date
	tomorrow := time.Now().Add(24 * time.Hour)
	times, err := i.iqamaForDate(tomorrow)
	handleErr(err)
	return &times, nil
}

func (i *IqamaCSV) GetDiscordPrettified() string {
	t := toTable(i.iqamaTimes)
	return "```markdown\n" + t.Render() + "\n```"
}

func (i *IqamaCSV) GetShellPrettified() string {
	t, err := i.GetTodayTimes()
	if err != nil {
		return fmt.Sprintf("Couldn't fetch today's iqama times %e", err)
	}
	return toTableDaily(*t).Render()
}

func (i *IqamaCSV) readCSV() error {
	file, err := os.Open(i.filePath)
	if err != nil {
		log.Fatalf("Unable to open file %s", i.filePath)
		return err
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("Unable to parse file as CSV for %s", i.filePath)
		return err
	}

	iqamaTimes := make(map[string]IqamaDailyTimes)
	for _, record := range records {
		if record[0] == "date" {
			continue
		}
		date, err := ParseDate(record[0])
		handleErr(err)
		fajrIqama, err := ParseHoursMinutes(record[1])
		handleErr(err)
		dhuhrIqama, err := ParseHoursMinutes(record[2])
		handleErr(err)
		asrIqama, err := ParseHoursMinutes(record[3])
		handleErr(err)
		maghribIqama, err := ParseHoursMinutes(record[4])
		handleErr(err)
		ishaIqama, err := ParseHoursMinutes(record[5])
		handleErr(err)

		dateStr := date.Format("01/02/2006")
		iqamaTimes[dateStr] = IqamaDailyTimes{
			Date: date,
			Fajr: Prayer{
				Iqama: fajrIqama,
			},
			Dhuhr: Prayer{
				Iqama: dhuhrIqama,
			},
			Asr: Prayer{
				Iqama: asrIqama,
			},
			Maghrib: Prayer{
				Iqama: maghribIqama,
			},
			Isha: Prayer{
				Iqama: ishaIqama,
			},
		}
	}

	i.iqamaTimes = iqamaTimes
	return nil
}

func handleErr(err error) {
	if err != nil {
		log.Fatalf("error: %s", err)
	}
}
