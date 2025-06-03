package main

import (
	"testing"
)

func TestBigInt(t *testing.T) {
	var userId int64 = 935247313479680099
	res := userId % 48

	t.Logf("res:%d", res)
}
