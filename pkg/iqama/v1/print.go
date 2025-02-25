package v1

import (
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
)

func GetDiscordPrettified() string {
	d, err := Get()
	if err != nil {
		return ""
	}
	t := toTable(*d)

	return "```markdown\n" + t.Render() + "\n```"
}

func GetShellPrettified() string {
	d, err := Get()
	if err != nil {
		return ""
	}
	t := toTable(*d)

	return t.Render()
}

func toTable(r Resp) table.Writer {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Prayer", "Adhan", "Iqama"})
	t.AppendRows([]table.Row{
		{"Fajr", r.Fajr.Adhan, r.Fajr.Iqama},
		{"Dhuhur", r.Dhuhr.Adhan, r.Dhuhr.Iqama},
		{"Asr", r.Asr.Adhan, r.Asr.Iqama},
		{"Maghrib", r.Maghrib.Adhan, r.Maghrib.Iqama},
		{"Isha", r.Isha.Adhan, r.Isha.Iqama},
	})
	return t
}
