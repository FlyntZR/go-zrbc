package utils

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
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

// Strtotime converts a date/time string to Unix timestamp
// It supports common formats like:
// - "2006-01-02 15:04:05"
// - "2006-01-02"
// - "2006/01/02 15:04:05"
// - "2006/01/02"
// - "now"
// - "+1 day"
// - "-1 day"
// - "+1 week"
// - "-1 week"
// - "+1 month"
// - "-1 month"
// - "+1 year"
// - "-1 year"
func Strtotime(str string) (int64, error) {
	if str == "" {
		return 0, errors.New("empty time string")
	}

	str = strings.TrimSpace(strings.ToLower(str))

	// Handle "now"
	if str == "now" {
		return time.Now().Unix(), nil
	}

	// Handle relative time strings
	if strings.HasPrefix(str, "+") || strings.HasPrefix(str, "-") {
		return parseRelativeTime(str)
	}

	// Try parsing common date formats
	layouts := []string{
		"2006-01-02 15:04:05",
		"2006-01-02",
		"2006/01/02 15:04:05",
		"2006/01/02",
		time.RFC3339,
		time.RFC1123Z,
		time.RFC1123,
		time.RFC822Z,
		time.RFC822,
		time.RubyDate,
		time.UnixDate,
		time.ANSIC,
	}

	for _, layout := range layouts {
		if t, err := time.ParseInLocation(layout, str, time.Local); err == nil {
			return t.Unix(), nil
		}
	}

	return 0, errors.New("unsupported time format")
}

// parseRelativeTime parses relative time strings like "+1 day", "-1 week" etc.
func parseRelativeTime(str string) (int64, error) {
	// Regular expression to match the pattern
	re := regexp.MustCompile(`^([+-])(\d+)\s+(day|week|month|year)s?$`)
	matches := re.FindStringSubmatch(str)

	if len(matches) != 4 {
		return 0, errors.New("invalid relative time format")
	}

	// Parse the number
	num, err := strconv.Atoi(matches[2])
	if err != nil {
		return 0, err
	}

	// If it's a negative offset, make the number negative
	if matches[1] == "-" {
		num = -num
	}

	now := time.Now()
	var t time.Time

	// Calculate the new time based on the unit
	switch matches[3] {
	case "day":
		t = now.AddDate(0, 0, num)
	case "week":
		t = now.AddDate(0, 0, num*7)
	case "month":
		t = now.AddDate(0, num, 0)
	case "year":
		t = now.AddDate(num, 0, 0)
	default:
		return 0, errors.New("unsupported time unit")
	}

	return t.Unix(), nil
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
