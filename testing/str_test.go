package main

import (
	"go-zrbc/pkg/utils"
	"testing"
)

func TestIsContainsSpecial(t *testing.T) {
	result := utils.IsContainsSpecial("'骗骗,喜欢你")
	// result := util.RemoveSpecail("Hello!@# World123")
	t.Logf("res:%v", result)
}

func TestRemoveSpecail(t *testing.T) {
	result := utils.RemoveSpecail("'骗骗?，：（,喜欢你")
	// result := util.RemoveSpecail("Hello!@# World123")
	t.Logf("res:%s", result)
}

func TestIsNumeric(t *testing.T) {
	testCases := []struct {
		input    string
		expected bool
	}{
		{"123", true},
		{"123.456", true},
		{"abc", false},
		{"123abc", false},
		{"-123.45", true},
		{"", false},
	}

	for _, tc := range testCases {
		result := utils.IsNumeric(tc.input)
		if result != tc.expected {
			t.Errorf("IsNumeric(%q) = %v; want %v", tc.input, result, tc.expected)
		}
	}
}
