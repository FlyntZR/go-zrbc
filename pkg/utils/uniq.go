package utils

import "regexp"

func UniqInt64(arr []int64) []int64 {
	var m = map[int64]struct{}{}
	for i := range arr {
		m[arr[i]] = struct{}{}
	}
	arr = arr[:0]
	for k := range m {
		arr = append(arr, k)
	}
	return arr
}

func IsMobile(mobile string) bool {
	result, _ := regexp.MatchString(`^(1[3|4|5|6|7|8|9][0-9]\d{4,8})$`, mobile)
	return result
}
