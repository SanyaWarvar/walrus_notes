package util

import "testing"

func TestGenerateRandomString(t *testing.T) {
	t.Run("correct length", func(t *testing.T) {
		for _, length := range []int{0, 1, 8, 32, 128} {
			got := GenerateRandomString(length)
			if len(got) != length {
				t.Errorf("GenerateRandomString(%d) length = %d, want %d", length, len(got), length)
			}
		}
	})

	t.Run("charset only", func(t *testing.T) {
		allowed := map[byte]bool{}
		for i := 0; i < len(charset); i++ {
			allowed[charset[i]] = true
		}

		got := GenerateRandomString(100)
		for i := 0; i < len(got); i++ {
			if !allowed[got[i]] {
				t.Errorf("GenerateRandomString() contains invalid char %q at pos %d", got[i], i)
			}
		}
	})

	t.Run("generates different values", func(t *testing.T) {
		a := GenerateRandomString(32)
		b := GenerateRandomString(32)
		if a == b {
			t.Error("expected two consecutive calls to produce different strings")
		}
	})
}
