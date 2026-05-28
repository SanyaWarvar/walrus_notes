package util

import (
	"net/http"
	"testing"
	"wn/pkg/constants"
)

func TestMaskBySize(t *testing.T) {
	tests := []struct {
		name string
		val  string
		want string
	}{
		{"empty", "", ""},
		{"single char", "a", "***"},
		{"three chars", "abc", "***"},
		{"four chars", "abcd", "abc*"},
		{"five chars", "abcde", "abc**"},
		{"nine chars", "123456789", "1234*****"},
		{"fifteen chars", "123456789012345", "12345**********"},
		{"sixteen chars", "1234567890123456", "123***456"},
		{"long token", "abcdefghijklmnopqrstuvwxyz", "abc***xyz"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := maskBySize(tt.val); got != tt.want {
				t.Errorf("maskBySize(%q) = %q, want %q", tt.val, got, tt.want)
			}
		})
	}
}

func TestMaskHeaders(t *testing.T) {
	headers := http.Header{}
	headers.Set("Content-Type", "application/json")
	headers.Set(constants.AuthorizationHeader, "Bearer secret-token-12345")
	headers.Set(constants.RefreshHeader, "refresh-token-value")

	masked := MaskHeaders(headers)

	if masked.Get("Content-Type") != "application/json" {
		t.Errorf("Content-Type should not be masked, got %q", masked.Get("Content-Type"))
	}

	authMasked := masked.Get(constants.AuthorizationHeader)
	if authMasked == headers.Get(constants.AuthorizationHeader) {
		t.Error("Authorization header should be masked")
	}
	if len(authMasked) == 0 {
		t.Error("masked Authorization should not be empty")
	}

	refreshMasked := masked.Get(constants.RefreshHeader)
	if refreshMasked == "" {
		t.Fatal("Refresh header should be present in masked headers")
	}
	if refreshMasked == headers.Get(constants.RefreshHeader) {
		t.Errorf("Refresh header should be masked, got %q", refreshMasked)
	}

	// Original headers must remain unchanged.
	if headers.Get(constants.AuthorizationHeader) != "Bearer secret-token-12345" {
		t.Error("MaskHeaders should not modify the original header map")
	}
}

func TestMaskHeaders_EmptySensitiveHeaders(t *testing.T) {
	headers := http.Header{
		"Content-Type": {"text/plain"},
	}

	masked := MaskHeaders(headers)
	if masked.Get("Content-Type") != "text/plain" {
		t.Errorf("unexpected Content-Type: %q", masked.Get("Content-Type"))
	}
}
