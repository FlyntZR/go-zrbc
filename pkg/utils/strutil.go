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
