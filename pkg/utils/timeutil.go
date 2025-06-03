package utils

import "time"

func TodayBegin() int64 {
	t := time.Now()
	addTime := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	return addTime.UnixMilli()
}

func WeekBegin() int64 {
	t := TodayBegin() - 7*24*3600000
	return t
}

func CurrentTime() (time.Time, error) {
	now := time.Now()
	shanghaiLocation, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return time.Time{}, err
	}
	return now.In(shanghaiLocation), nil
}
