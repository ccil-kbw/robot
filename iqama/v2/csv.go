package v2

import (
	"encoding/csv"
	"log"
	"os"
	"time"
)

type IqamaCSV struct {
	filePath   string
	iqamaTimes map[time.Time]IqamaDailyTimes
}

func NewIqamaCSV(filePath string) Iqama {
	i := &IqamaCSV{filePath: filePath}
	if err := i.readCSV(); err != nil {
		log.Fatalf("Unable to read CSV file %s", filePath)
	}
	return i
}

func (i *IqamaCSV) GetTodayTimes() (*IqamaDailyTimes, error) {
	// Get today's date
	today := time.Now().Truncate(24 * time.Hour)
	// Get today's iqama times
	times, ok := i.iqamaTimes[today]
	if !ok {
		return nil, nil
	}
	return &times, nil
}

func (i *IqamaCSV) GetTomorrowTimes() (*IqamaDailyTimes, error) {
	// Get tomorrow's date
	tomorrow := time.Now().Add(24 * time.Hour).Truncate(24 * time.Hour)
	// Get tomorrow's iqama times
	times, ok := i.iqamaTimes[tomorrow]
	if !ok {
		return nil, nil
	}
	return &times, nil
}

func (i *IqamaCSV) GetDiscordPrettified() string {
	t := toTable(i.iqamaTimes)
	return "```markdown\n" + t.Render() + "\n```"
}

func (i *IqamaCSV) GetShellPrettified() string {
	//TODO implement me
	panic("implement me")
}

func (i *IqamaCSV) readCSV() error {
	file, err := os.Open(i.filePath)
	if err != nil {
		log.Fatalf("Unable to open file %s", i.filePath)
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("Unable to parse file as CSV for %s", i.filePath)
		return err
	}

	iqamaTimes := make(map[time.Time]IqamaDailyTimes)
	for _, record := range records {
		// Convert each record to IqamaDailyTimes and add to iqamaTimes
		// You need to replace the following lines with your conversion logic
		date, _ := time.Parse("2006-01-02", record[0]) // assuming the date is the first field in the record
		iqamaTimes[date] = IqamaDailyTimes{}
	}

	i.iqamaTimes = iqamaTimes

	return nil
}
