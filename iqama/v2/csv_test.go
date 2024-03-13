package v2

import (
	"reflect"
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
	type fields struct {
		filePath   string
		iqamaTimes map[time.Time]IqamaDailyTimes
	}
	tests := []struct {
		name    string
		fields  fields
		want    *IqamaDailyTimes
		wantErr bool
	}{
		{
			name: "TestIqamaCSV_GetTodayTimes",
			fields: fields{
				filePath:   "test_assets/iqama_2024.csv",
				iqamaTimes: nil,
			},
			want:    nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &IqamaCSV{
				filePath:   tt.fields.filePath,
				iqamaTimes: tt.fields.iqamaTimes,
			}
			got, err := i.GetTodayTimes()
			if (err != nil) != tt.wantErr {
				t.Errorf("IqamaCSV.GetTodayTimes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IqamaCSV.GetTodayTimes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIqamaCSV_GetTomorrowTimes(t *testing.T) {
	type fields struct {
		filePath   string
		iqamaTimes map[time.Time]IqamaDailyTimes
	}
	tests := []struct {
		name    string
		fields  fields
		want    *IqamaDailyTimes
		wantErr bool
	}{
		{
			name: "TestIqamaCSV_GetTomorrowTimes",
			fields: fields{
				filePath:   "test_assets/iqama_2024.csv",
				iqamaTimes: nil,
			},
			want:    nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &IqamaCSV{
				filePath:   tt.fields.filePath,
				iqamaTimes: tt.fields.iqamaTimes,
			}
			got, err := i.GetTomorrowTimes()
			if (err != nil) != tt.wantErr {
				t.Errorf("IqamaCSV.GetTomorrowTimes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IqamaCSV.GetTomorrowTimes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIqamaCSV_GetDiscordPrettified(t *testing.T) {
	type fields struct {
		filePath   string
		iqamaTimes map[time.Time]IqamaDailyTimes
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

	x := time.RFC3339
}
