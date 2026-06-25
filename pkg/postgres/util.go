package postgres

import (
	"time"

	"github.com/jackc/pgx/v5"
)

func nullString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func nullTime(t time.Time) *time.Time {
	if t.IsZero() {
		return nil
	}
	return &t
}

func isNoRows(err error) bool {
	return err == pgx.ErrNoRows
}
