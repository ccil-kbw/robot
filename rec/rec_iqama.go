package rec

import (
	"fmt"
	"time"

	v1 "github.com/ccil-kbw/robot/iqama/v1"
)

var (
	EveryDay             []time.Weekday = []time.Weekday{time.Sunday, time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday, time.Saturday}
	JumuaaRecordDuration time.Duration  = 2 * time.Hour
	DarsRecordDuration   time.Duration  = 1 * time.Hour
	location             string         = "America/Montreal"
)

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
func SupposedToBeRecording(confs []*RecordConfig) bool {

	// Added this outside of the for loop to have better logging internally by looking at all records before returning
	shouldRecord := false

	now := time.Now()

	fmt.Printf("current time: %s \n", now.Format("15:04:05"))
	for _, conf := range confs {
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

		fmt.Printf("%s %v %v ", conf.Description, conf.StartTime.Format("15:04:05"), conf.StartTime.Add(conf.Duration).Format("15:04:05"))

		// Check if we're in the time range (from conf.StartTime to conf.StartTime+Duration)
		if now.After(startTime) && now.Before(startTime.Add(conf.Duration)) {
			shouldRecord = true
		}

		fmt.Printf("--- in time range: %v\n", shouldRecord)

	}

	return shouldRecord
}

func GetIqamaRecordingConfigs() []*RecordConfig {
	timeLocation, err := time.LoadLocation("America/Montreal")
	if err != nil {
		panic(err)
	}

	iqamaTimes := v1.Get()

	fmt.Printf("fajr time: %s\n", iqamaTimes.Fajr.Iqama)
	fmt.Printf("dhuhur time: %s\n", iqamaTimes.Dhuhr.Iqama)
	fmt.Printf("isha time: %s\n", iqamaTimes.Isha.Iqama)

	return []*RecordConfig{
		{
			Description:   "Fajr Recording",
			StartTime:     toTime(iqamaTimes.Fajr.Iqama),
			Duration:      DarsRecordDuration,
			RecordingDays: EveryDay,
		},
		{
			Description:   "Dhuhur Recording",
			StartTime:     toTime(iqamaTimes.Dhuhr.Iqama),
			Duration:      DarsRecordDuration,
			RecordingDays: EveryDay,
		},
		{
			Description:   "Isha Recording",
			StartTime:     toTime(iqamaTimes.Isha.Iqama),
			Duration:      DarsRecordDuration,
			RecordingDays: EveryDay,
		},
		{
			Description:   "Jumuaa Recording",
			StartTime:     time.Date(2023, 1, 1, 11, 55, 0, 0, timeLocation),
			Duration:      JumuaaRecordDuration,
			RecordingDays: []time.Weekday{time.Friday},
		},
		{
			Description:   "Fiqh Dars Recording",
			StartTime:     toTime(iqamaTimes.Maghrib.Iqama),
			Duration:      JumuaaRecordDuration,
			RecordingDays: []time.Weekday{time.Thursday},
		},
	}
}

func toTime(s string) time.Time {
	t, err := time.Parse("15:04", s)
	if err != nil {
		panic(fmt.Sprintf("could not parse prayer time %s, %s", s, err.Error()))
	}
	fmt.Println(timeToday(t.Hour(), t.Minute()))
	return timeToday(t.Hour(), t.Minute())
}

func timeToday(hour, minute int) time.Time {
	now := time.Now()

	l, err := time.LoadLocation(location)
	if err != nil {
		panic(err)
	}
	return time.Date(now.Year(), now.Month(), now.Day(), hour, minute, 0, 0, l)
}
