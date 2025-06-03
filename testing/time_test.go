package main

import (
	"testing"
	"time"
)

func TestTime(t *testing.T) {
	cTime := time.UnixMilli(1709538837 * 1000)

	ts := cTime.UnixMilli()
	t.Logf("%d", ts)
}
