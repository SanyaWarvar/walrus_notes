package util

import (
	"fmt"
	"time"
)

func GetCurrentMskTime() time.Time {
	return time.Now().UTC().Add(3 * time.Hour)
}

func GetCurrentUTCTime() time.Time {
	return time.Now().UTC()
}

func ConvertStringToTime(timeStr string) (*time.Time, error) {
	t, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return nil, fmt.Errorf("cant convert string {%s} to time", timeStr)
	}
	return &t, nil
}
