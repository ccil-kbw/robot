package v2

import (
	"fmt"
	"os"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
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

// GetDefaultIqamaCSVPath returns the path to the iqama CSV file
// Priority: 1. MDROID_IQAMA_CSV_PATH environment variable
//           2. iqama_YYYY.csv where YYYY is current year
func GetDefaultIqamaCSVPath() string {
	// Check environment variable first
	if path := os.Getenv("MDROID_IQAMA_CSV_PATH"); path != "" {
		return path
	}

	// Default to current year
	year := time.Now().Year()
	return fmt.Sprintf("iqama_%d.csv", year)
}
