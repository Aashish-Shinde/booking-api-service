package utils

import (
	"testing"
	"time"
)

func TestTimeStringToTime(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid morning time", "10:30", false},
		{"valid afternoon time", "14:45", false},
		{"invalid format", "10-30", true},
		{"invalid hour", "25:00", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := TimeStringToTime(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("TimeStringToTime() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCheckTimeAlignment(t *testing.T) {
	tests := []struct {
		name string
		time string
		want bool
	}{
		{"aligned to 30 min 1", "10:00", true},
		{"aligned to 30 min 2", "10:30", true},
		{"not aligned 1", "14:45", false},
		{"not aligned 2", "10:15", false},
		{"not aligned 3", "10:45", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CheckTimeAlignment(tt.time)
			if got != tt.want {
				t.Errorf("CheckTimeAlignment(%q) = %v, want %v", tt.time, got, tt.want)
			}
		})
	}
}

func TestTimeStrToMinutes(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    int
		wantErr bool
	}{
		{"midnight", "00:00", 0, false},
		{"morning", "10:30", 630, false},
		{"afternoon", "14:45", 885, false},
		{"end of day", "23:59", 1439, false},
		{"invalid format", "10-30", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := TimeStrToMinutes(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("TimeStrToMinutes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("TimeStrToMinutes(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestMinutesToTimeStr(t *testing.T) {
	tests := []struct {
		name string
		mins int
		want string
	}{
		{"midnight", 0, "00:00"},
		{"morning", 630, "10:30"},
		{"afternoon", 885, "14:45"},
		{"end of day", 1439, "23:59"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MinutesToTimeStr(tt.mins)
			if got != tt.want {
				t.Errorf("MinutesToTimeStr(%d) = %s, want %s", tt.mins, got, tt.want)
			}
		})
	}
}

func TestGetDayOfWeek(t *testing.T) {
	// April 7, 2026 is a Tuesday
	date := time.Date(2026, time.April, 7, 0, 0, 0, 0, time.UTC)
	expected := 2

	got := GetDayOfWeek(date)
	if got != expected {
		t.Errorf("GetDayOfWeek() = %v, want %v", got, expected)
	}
}

func TestAddMinutes(t *testing.T) {
	base := time.Date(2026, time.April, 7, 10, 0, 0, 0, time.UTC)
	result := AddMinutes(base, 30)

	expected := time.Date(2026, time.April, 7, 10, 30, 0, 0, time.UTC)
	if !result.Equal(expected) {
		t.Errorf("AddMinutes() = %v, want %v", result, expected)
	}
}
