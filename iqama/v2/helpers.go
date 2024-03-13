package v2

import (
	"github.com/jedib0t/go-pretty/v6/table"
	"os"
	"time"
)

func toTable(m map[time.Time]IqamaDailyTimes) table.Writer {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Date", "Fajr", "Dhuhur", "Asr", "Maghrib", "Isha"})
	for date, times := range m {
		t.AppendRows([]table.Row{
			{date.Format("2006-01-02"), times.Fajr, times.Dhuhr, times.Asr, times.Maghrib, times.Isha},
		})
	}
	return t
}
