package v2

import (
	"github.com/jedib0t/go-pretty/v6/table"
	"os"
	"time"
)

func toTable(m map[string]IqamaDailyTimes) table.Writer {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Date", "Fajr", "Dhuhur", "Asr", "Maghrib", "Isha"})
	for date, times := range m {
		t.AppendRows([]table.Row{
			{date, times.Fajr, times.Dhuhr, times.Asr, times.Maghrib, times.Isha},
		})
	}
	return t
}

func toTableDaily(m IqamaDailyTimes) table.Writer {
	t := table.NewWriter()
	t.AppendHeader(table.Row{"Date", "Fajr", "Dhuhur", "Asr", "Maghrib", "Isha"})
	t.AppendRows([]table.Row{
		{FormatDate(m.Date), FormatTime(m.Fajr.Iqama), FormatTime(m.Dhuhr.Iqama), FormatTime(m.Asr.Iqama), FormatTime(m.Maghrib.Iqama), FormatTime(m.Isha.Iqama)}})
	return t
}

// ParseHoursMinutes parses a string in the format "15:04" and returns a time.Time
func ParseHoursMinutes(s string) (time.Time, error) {
	return time.Parse("3:04 pm", s)
}

// ParseDate parses a string in the format "01/02/2006" (Month/Day/Year) and returns a time.Time
func ParseDate(s string) (time.Time, error) {
	return time.Parse("01/02/2006", s)
}

func FormatDate(t time.Time) string {
	return t.Format("01/02/2006")
}

func FormatTime(t time.Time) string {
	return t.Format("3:04 pm")
}
