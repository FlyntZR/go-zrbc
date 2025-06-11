package main

import (
	"go-zrbc/pkg/utils"
	"testing"
	"time"
)

func TestTime(t *testing.T) {
	cTime := time.UnixMilli(1709538837 * 1000)

	ts := cTime.UnixMilli()
	t.Logf("%d", ts)
}

func TestStrtotime(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "standard datetime format",
			input:   "2024-03-04 15:04:05",
			wantErr: false,
		},
		{
			name:    "date only format",
			input:   "2024-03-04",
			wantErr: false,
		},
		{
			name:    "now keyword",
			input:   "now",
			wantErr: false,
		},
		{
			name:    "relative time future",
			input:   "+1 day",
			wantErr: false,
		},
		{
			name:    "relative time past",
			input:   "-1 week",
			wantErr: false,
		},
		{
			name:    "empty string",
			input:   "",
			wantErr: true,
		},
		{
			name:    "invalid format",
			input:   "invalid-time-format",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := utils.Strtotime(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Strtotime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got <= 0 {
				t.Errorf("Strtotime() got invalid timestamp = %v", got)
			}
			t.Logf("Strtotime() got = %v", got)
		})
	}
}
