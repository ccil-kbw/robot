package v2

import "time"

type IqamaDailyTimes struct {
	Date    time.Time
	Fajr    Prayer
	Sunrise time.Time
	Dhuhr   Prayer
	Asr     Prayer
	Maghrib Prayer
	Isha    Prayer
	Jumuaa  Jumuaa
}

type Prayer struct {
	Adhan time.Time
	Iqama time.Time
}

type Jumuaa struct {
	Times []time.Time
}
