package common

import (
	"errors"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
)

func TestIsUniqueErr(t *testing.T) {
	t.Run("unique violation", func(t *testing.T) {
		err := &pgconn.PgError{Code: "23505"}
		if !IsUniqueErr(err) {
			t.Error("IsUniqueErr() = false, want true")
		}
	})

	t.Run("other pg error", func(t *testing.T) {
		err := &pgconn.PgError{Code: "23503"}
		if IsUniqueErr(err) {
			t.Error("IsUniqueErr() = true, want false")
		}
	})

	t.Run("non pg error", func(t *testing.T) {
		if IsUniqueErr(errors.New("generic error")) {
			t.Error("IsUniqueErr() = true, want false")
		}
	})
}
