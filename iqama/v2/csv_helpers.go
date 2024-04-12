package v2

import (
	"errors"
	"time"
)

func (i *IqamaCSV) iqamaForDate(date time.Time) (IqamaDailyTimes, error) {
	times, ok := i.iqamaTimes[date.Format("01/02/2006")]
	if !ok {
		return IqamaDailyTimes{}, errors.New("could not get iqama times for date")
	}
	return times, nil
}
