package iqama

import (
	"fmt"
	v1 "github.com/ccil-kbw/robot/iqama/v1"
	"github.com/ccil-kbw/robot/rec"
	"sync"
	"time"
)

// RecordingsData (sync.Mutex) is used to handle the RecordConfig asynchronously as it's used in the main and some sub routines
type RecordingsData struct {
	mu    sync.Mutex
	confs *[]rec.RecordConfig
}

// Refresh updates the scheduling configurations
func (ac *RecordingsData) Refresh() {
	ac.mu.Lock()
	ac.confs = rec.RecordConfigData.Get()
	ac.mu.Unlock()
}

func (ac *RecordingsData) Confs() *[]rec.RecordConfig {
	ac.mu.Lock()
	defer ac.mu.Unlock()
	return ac.confs
}

// PrayersData (sync.Mutex) is used to handle the Iqama times
type PrayersData struct {
	mu    sync.Mutex
	confs *v1.Resp
}

func (i *PrayersData) Refresh() {
	i.mu.Lock()
	defer i.mu.Unlock()
	resp, _ := v1.Get()
	i.confs = resp
}

func (i *PrayersData) Confs() *v1.Resp {
	i.mu.Lock()
	defer i.mu.Unlock()
	return i.confs
}

func StartRecordingScheduleServer() *RecordingsData {
	data := &RecordingsData{}
	data.Refresh()
	ticker := time.NewTicker(1 * time.Hour)
	go func() {
		for range ticker.C {
			fmt.Printf("%s: Updating Iqama time...\n", time.Now().Format(time.Kitchen))
			data.Refresh()
		}
	}()
	return data
}

func StartIqamaServer() *PrayersData {
	data := &PrayersData{}
	data.Refresh()
	ticker := time.NewTicker(1 * time.Hour)
	go func() {
		for range ticker.C {
			fmt.Printf("%s: Updating Iqama time...\n", time.Now().Format(time.Kitchen))
			data.Refresh()
		}
	}()
	return data
}
