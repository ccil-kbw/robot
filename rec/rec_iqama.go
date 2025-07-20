package rec

import (
	"fmt"
	v2 "github.com/ccil-kbw/robot/iqama/v2"
	"sync"
	"time"
)

var (
	EveryDayButFriday    []time.Weekday = []time.Weekday{time.Sunday, time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Saturday}
	JumuaaRecordDuration time.Duration  = 2 * time.Hour
	DarsRecordDuration   time.Duration  = 35 * time.Minute
	location             string         = "America/Montreal"
)

type RecordConfigDataS struct {
	data  *[]RecordConfig
	iqama v2.Iqama
	mu    sync.Mutex
}

// NewRecordConfigDataS TODO: temporarily in code but needs to move to the postgresql db from the robot-api project
func NewRecordConfigDataS() *RecordConfigDataS {
	iqamaClient := v2.NewIqamaCSV("iqama_2025.csv")
	today, err := iqamaClient.GetTodayTimes()
	if err != nil {
		fmt.Println("couldn't fetch iqama times, keeping current data")
		return nil
	}

	if today == nil {
		fmt.Println("couldn't fetch iqama times, keeping current data")
		return nil
	}

	fajr := today.Fajr.Iqama
	dhuhr := today.Dhuhr.Iqama
	maghrib := today.Maghrib.Iqama
	isha := today.Isha.Iqama
	rc := &RecordConfigDataS{
		iqama: iqamaClient,
		data: &[]RecordConfig{
			{
				Description:   "Jumuaa Recording",
				StartTime:     time.Date(2024, 1, 1, 11, 55, 0, 0, time.Local),
				Duration:      JumuaaRecordDuration,
				RecordingDays: []time.Weekday{time.Friday},
			},
			{
				Description:   "Fajr Recording",
				StartTime:     fajr,
				Duration:      DarsRecordDuration,
				RecordingDays: EveryDayButFriday,
			},
			{
				Description:   "Dhuhur Recording",
				StartTime:     dhuhr,
				Duration:      DarsRecordDuration,
				RecordingDays: EveryDayButFriday,
			},
			{
				Description:   "Maghrib Recording Sunday",
				StartTime:     maghrib,
				Duration:      isha.Add(-5 * time.Minute).Sub(maghrib),
				RecordingDays: []time.Weekday{time.Saturday, time.Sunday},
			},
			{
				Description:   "Isha Recording",
				StartTime:     isha,
				Duration:      DarsRecordDuration,
				RecordingDays: EveryDayButFriday,
			},
		},
	}
	return rc
}

func (rc *RecordConfigDataS) Get() *[]RecordConfig {
	defer rc.mu.Unlock()
	rc.mu.Lock()
	return rc.data
}

type RecordConfig struct {
	Description string

	// Only check Hours and Minutes
	StartTime time.Time

	// Duration of the Recording,
	// StopTime is StartTime + Duration
	Duration time.Duration

	// For everyday use helper EveryDay)
	RecordingDays []time.Weekday
}

// SupposedToBeRecording just what the func name is saying.
// Please add doc wherever you think it was unreadable, else refactor the portion
func SupposedToBeRecording(data *RecordConfigDataS) bool {

	// Added this outside of the for loop to have better logging internally by looking at all records before returning
	shouldRecord := false

	now := time.Now()

	fmt.Printf("current time: %s \n", now.Format("15:04:05"))
	for _, conf := range *data.Get() {
		recordToday := false

		// Check if should be recording today
		for _, day := range conf.RecordingDays {
			if day == time.Now().Weekday() {
				recordToday = true
			}
		}

		if !recordToday {
			continue
		}

		// Set today's start time for this prayer time for Today
		startTime := timeToday(conf.StartTime.Hour(), conf.StartTime.Minute())

		fmt.Printf("%s %v %v \n", conf.Description, conf.StartTime.Format("15:04:05"), conf.StartTime.Add(conf.Duration).Format("15:04:05"))

		// Check if we're in the time range (from conf.StartTime to conf.StartTime+Duration)
		if now.After(startTime) && now.Before(startTime.Add(conf.Duration)) {
			shouldRecord = true
			break
		}

	}

	return shouldRecord
}

func timeToday(hour, minute int) time.Time {
	now := time.Now()

	l, err := time.LoadLocation(location)
	if err != nil {
		panic(err)
	}
	return time.Date(now.Year(), now.Month(), now.Day(), hour, minute, 0, 0, l)
}
