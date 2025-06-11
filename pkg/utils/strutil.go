package utils

import (
	"regexp"
	"strings"
)

var specialChars = []string{`,`, `，`, `：`, `（`, `）`, `(`, `)`, `:`, `?`, `？`, `*`, `|`, `{`, `}`, `.`, `\`, `/`, `$`, `、`, `·`, "`", `'`, `"`}

func IsContainsSpecial(orgStr string) bool {
	if strings.ContainsAny(orgStr, strings.Join(specialChars, "")) {
		return true
	} else {
		return false
	}
}

func RemoveSpecail(orgStr string) string {
	re := regexp.MustCompile("[,，：:（）()?？*|{}.\\/$、·`'\"]")
	return re.ReplaceAllString(orgStr, "")
}

func IsNumeric(str string) bool {
	// 匹配整数或浮点数
	matched, _ := regexp.MatchString(`^-?\d*\.?\d+$`, str)
	return matched
}
