package util

import (
	"testing"
	"time"
)

func TestConvertStringToTime(t *testing.T) {
	t.Run("valid RFC3339", func(t *testing.T) {
		input := "2024-06-15T12:30:00Z"
		got, err := ConvertStringToTime(input)
		if err != nil {
			t.Fatalf("ConvertStringToTime() error = %v", err)
		}
		want, _ := time.Parse(time.RFC3339, input)
		if !got.Equal(want) {
			t.Errorf("ConvertStringToTime() = %v, want %v", got, want)
		}
	})

	t.Run("invalid format", func(t *testing.T) {
		_, err := ConvertStringToTime("not-a-date")
		if err == nil {
			t.Fatal("expected error for invalid time string")
		}
	})
}

func TestGetCurrentUTCTime(t *testing.T) {
	before := time.Now().UTC()
	got := GetCurrentUTCTime()
	after := time.Now().UTC()

	if got.Before(before) || got.After(after) {
		t.Errorf("GetCurrentUTCTime() = %v, expected between %v and %v", got, before, after)
	}
}

func TestGetCurrentMskTime(t *testing.T) {
	utc := GetCurrentUTCTime()
	msk := GetCurrentMskTime()

	diff := msk.Sub(utc)
	want := 3 * time.Hour
	if diff < want-time.Second || diff > want+time.Second {
		t.Errorf("MSK offset = %v, want ~%v", diff, want)
	}
}
