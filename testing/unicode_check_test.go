package main

import (
	"testing"
	"unicode"
)

func TestUnicodeCheck(t *testing.T) {
	v := "恋爱麻烦但。是甜"
	r := []rune(v)
	for i := 0; i < len(r); i++ {
		if unicode.IsPunct(r[i]) {
			t.Logf("true, %v", string(r[:i]))
		} else {
			t.Log("false")
		}
	}
}

func TestStringCheck(t *testing.T) {
	v := "abcde"
	t.Logf("v: %s", v[:3])
}
