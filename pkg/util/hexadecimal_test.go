package util

import "testing"

func TestHexadecimalWithPadding(t *testing.T) {
	tests := []struct {
		name string
		n    int
		want string
	}{
		{"zero padded to 6", 0, "#####0"},
		{"single digit", 1, "#####1"},
		{"two hex digits padded", 255, "####ff"},
		{"exactly 6 hex chars", 0xffffff, "ffffff"},
		{"more than 6 hex chars", 0x1000000, "1000000"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HexadecimalWithPadding(tt.n); got != tt.want {
				t.Errorf("HexadecimalWithPadding(%d) = %q, want %q", tt.n, got, tt.want)
			}
		})
	}
}

func TestIsHexNumber(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want bool
	}{
		{"valid lowercase", "deadbeef", true},
		{"valid uppercase", "DEADBEEF", true},
		{"valid mixed", "DeAdBeEf", true},
		{"empty string", "", true},
		{"odd length", "abc", false},
		{"non-hex chars", "ghij", false},
		{"with spaces", "ab cd", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsHexNumber(tt.s); got != tt.want {
				t.Errorf("IsHexNumber(%q) = %v, want %v", tt.s, got, tt.want)
			}
		})
	}
}
