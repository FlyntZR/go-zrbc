package utils

import "strconv"

func ParseInt(bigIntStr string) (int64, error) {
	if bigIntStr == "" {
		return 0, nil
	}
	return strconv.ParseInt(bigIntStr, 10, 64)
}
