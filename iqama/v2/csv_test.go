package v2

import (
	"testing"
	"time"
)

func TestNewIqamaCSV(t *testing.T) {
	type args struct {
		filePath string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "TestNewIqamaCSV",
			args: args{
				filePath: "test_assets/iqama_2024.csv",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			NewIqamaCSV(tt.args.filePath)
		})
	}
}

func TestIqamaCSV_GetTodayTimes(t *testing.T) {
	// Manually Load CSV and read Today's Record
	i := NewIqamaCSV("test_assets/iqama_2024.csv")
	daily, err := i.GetTodayTimes()
	if err != nil {
		t.Errorf("IqamaCSV.GetTodayTimes() error = %v", err)
	}

	if daily.Date.Format(DateFormat) != time.Now().Format(DateFormat) {
		t.Errorf("IqamaCSV.GetTodayTimes() = %v, want %v", daily.Date, time.Now().Format("01/02/2006"))
	}

	// For the following tests we will only look at the hours of the time.Time objects
	// This will allow us to compare the times without worrying about the date

	// Assert that Fajr is after 3:00 am (Summer) and before 7am (Winter)
	if daily.Fajr.Adhan.Hour() < 3 || daily.Fajr.Adhan.Hour() > 6 {
		t.Errorf("IqamaCSV.GetTodayTimes() = %v, want between 3am and 6am", daily.Fajr.Adhan)
	}

	// Assert that Dhuhr is after 11:00 am and before 2pm
	if daily.Dhuhr.Adhan.Hour() < 11 || daily.Dhuhr.Adhan.Hour() > 14 {
		t.Errorf("IqamaCSV.GetTodayTimes() = %v, want between 12pm and 1pm", daily.Dhuhr.Adhan)
	}

	// Assert that Asr is after 2:00 pm and before 5pm
	if daily.Asr.Adhan.Hour() < 14 || daily.Asr.Adhan.Hour() > 17 {
		t.Errorf("IqamaCSV.GetTodayTimes() = %v, want between 2pm and 5pm", daily.Asr.Adhan)
	}

	// Assert that Maghrib is after 4:00 pm and before 9pm
	if daily.Maghrib.Adhan.Hour() < 16 || daily.Maghrib.Adhan.Hour() > 21 {
		t.Errorf("IqamaCSV.GetTodayTimes() = %v, want between 5pm and 8pm", daily.Maghrib.Adhan)
	}

	// Assert that Isha is after 5:00 pm and before 11pm
	if daily.Isha.Adhan.Hour() < 17 || daily.Isha.Adhan.Hour() > 23 {
		t.Errorf("IqamaCSV.GetTodayTimes() = %v, want between 6pm and 9pm", daily.Isha.Adhan)

	}
}

func TestIqamaCSV_GetDiscordPrettified(t *testing.T) {
	type fields struct {
		filePath   string
		iqamaTimes map[string]IqamaDailyTimes
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "TestIqamaCSV_GetDiscordPrettified",
			fields: fields{
				filePath:   "test_assets/iqama_2024.csv",
				iqamaTimes: nil,
			},
			want: "```markdown\n+------+------+--------+-----+---------+------+\n| DATE | FAJR | DHUHUR | ASR | MAGHRIB | ISHA |\n+------+------+--------+-----+---------+------+\n+------+------+--------+-----+---------+------+\n```",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &IqamaCSV{
				filePath:   tt.fields.filePath,
				iqamaTimes: tt.fields.iqamaTimes,
			}
			if got := i.GetDiscordPrettified(); got != tt.want {
				t.Errorf("IqamaCSV.GetDiscordPrettified() = %v, want %v", got, tt.want)
			}
		})
	}
}
