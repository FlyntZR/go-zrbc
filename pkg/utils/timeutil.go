package utils

import (
	"time"
)

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

// CheckTimestamp validates a Unix timestamp ensuring it's:
// 1. Not empty
// 2. A valid 10-digit Unix timestamp
// 3. Within acceptable time range (-4 to 30 seconds from current time)
func CheckTimestamp(timestamp int64) error {
	// Get current timestamp
	now := time.Now().Unix()

	// If timestamp is 0, return without error
	if timestamp == 0 {
		return ErrTimeFormatError
	}

	// Check if timestamp is valid (10 digits)
	if timestamp < 1000000000 || timestamp > 9999999999 {
		return ErrTimestampError
	}

	// Check time difference
	diff := now - timestamp
	if diff < -4 || diff > 30 {
		return ErrTimestampError
	}

	return nil
}
