package util

import (
	"testing"
	"wn/pkg/constants"
)

func TestCalculateOffset(t *testing.T) {
	tests := []struct {
		name string
		page int
		want int
	}{
		{"first page", 1, 0},
		{"second page", 2, constants.PageSize},
		{"third page", 3, 2 * constants.PageSize},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CalculateOffset(tt.page); got != tt.want {
				t.Errorf("CalculateOffset(%d) = %d, want %d", tt.page, got, tt.want)
			}
		})
	}
}

func TestCalculateLimit(t *testing.T) {
	if got := CalculateLimit(); got != constants.PageSize {
		t.Errorf("CalculateLimit() = %d, want %d", got, constants.PageSize)
	}
}
