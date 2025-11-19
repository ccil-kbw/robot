package v2

import (
	"encoding/csv"
	"fmt"
	"os"
	"time"
)

// IqamaCSV implements the Iqama interface by reading prayer times from a CSV file.
// The CSV file should contain columns: date, adhanFajr, iqamaFajr, shourouk, adhanDhuhr,
// iqamaDhuhr, adhanAsr, iqamaAsr, adhanMaghrib, iqamaMaghrib, adhanIsha, iqamaIsha.
type IqamaCSV struct {
	filePath   string
	iqamaTimes map[string]IqamaDailyTimes
}

// NewIqamaCSV creates a new IqamaCSV instance by loading prayer times from the specified CSV file.
// It reads and parses the entire CSV file during initialization.
// Returns an error if the file cannot be read or parsed.
func NewIqamaCSV(filePath string) (Iqama, error) {
	i := &IqamaCSV{filePath: filePath}
	if err := i.readCSV(); err != nil {
		return nil, fmt.Errorf("unable to read CSV file %s: %w", filePath, err)
	}
	return i, nil
}

// GetTodayTimes returns the prayer times for the current day.
// Returns an error if today's date is not found in the CSV or if parsing fails.
func (i *IqamaCSV) GetTodayTimes() (*IqamaDailyTimes, error) {
	today := time.Now()
	times, err := i.iqamaForDate(today)
	if err != nil {
		return nil, fmt.Errorf("failed to get today's iqama times: %w", err)
	}
	return &times, nil
}

// GetTomorrowTimes returns the prayer times for tomorrow.
// Returns an error if tomorrow's date is not found in the CSV or if parsing fails.
func (i *IqamaCSV) GetTomorrowTimes() (*IqamaDailyTimes, error) {
	tomorrow := time.Now().Add(24 * time.Hour)
	times, err := i.iqamaForDate(tomorrow)
	if err != nil {
		return nil, fmt.Errorf("failed to get tomorrow's iqama times: %w", err)
	}
	return &times, nil
}

// GetDiscordPrettified returns all prayer times formatted as a Discord-friendly markdown table.
// The output is wrapped in markdown code blocks for proper Discord rendering.
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
		return fmt.Errorf("unable to open file %s: %w", i.filePath, err)
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("unable to parse file as CSV for %s: %w", i.filePath, err)
	}

	iqamaTimes := make(map[string]IqamaDailyTimes)
	for _, record := range records {
		if record[0] == "date" {
			continue
		}
		date, err := ParseDate(record[0])
		if err != nil {
			return fmt.Errorf("failed to parse date: %w", err)
		}
		fajrAdhan, err := ParseHoursMinutes(record[1])
		if err != nil {
			return fmt.Errorf("failed to parse Fajr Adhan: %w", err)
		}
		fajrIqama, err := ParseHoursMinutes(record[2])
		if err != nil {
			return fmt.Errorf("failed to parse Fajr Iqama: %w", err)
		}
		dhuhrAdhan, err := ParseHoursMinutes(record[4])
		if err != nil {
			return fmt.Errorf("failed to parse Dhuhr Adhan: %w", err)
		}
		dhuhrIqama, err := ParseHoursMinutes(record[5])
		if err != nil {
			return fmt.Errorf("failed to parse Dhuhr Iqama: %w", err)
		}
		asrAdhan, err := ParseHoursMinutes(record[6])
		if err != nil {
			return fmt.Errorf("failed to parse Asr Adhan: %w", err)
		}
		asrIqama, err := ParseHoursMinutes(record[7])
		if err != nil {
			return fmt.Errorf("failed to parse Asr Iqama: %w", err)
		}
		maghribAdhan, err := ParseHoursMinutes(record[8])
		if err != nil {
			return fmt.Errorf("failed to parse Maghrib Adhan: %w", err)
		}
		maghribIqama, err := ParseHoursMinutes(record[9])
		if err != nil {
			return fmt.Errorf("failed to parse Maghrib Iqama: %w", err)
		}
		ishaAdhan, err := ParseHoursMinutes(record[10])
		if err != nil {
			return fmt.Errorf("failed to parse Isha Adhan: %w", err)
		}
		ishaIqama, err := ParseHoursMinutes(record[11])
		if err != nil {
			return fmt.Errorf("failed to parse Isha Iqama: %w", err)
		}

		dateStr := date.Format("01/02/2006")
		iqamaTimes[dateStr] = IqamaDailyTimes{
			Date: date,
			Fajr: Prayer{
				Adhan: fajrAdhan,
				Iqama: fajrIqama,
			},
			Dhuhr: Prayer{
				Adhan: dhuhrAdhan,
				Iqama: dhuhrIqama,
			},
			Asr: Prayer{
				Adhan: asrAdhan,
				Iqama: asrIqama,
			},
			Maghrib: Prayer{
				Adhan: maghribAdhan,
				Iqama: maghribIqama,
			},
			Isha: Prayer{
				Adhan: ishaAdhan,
				Iqama: ishaIqama,
			},
		}
	}

	i.iqamaTimes = iqamaTimes
	return nil
}
