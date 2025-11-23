package main

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func NumericToString(n pgtype.Numeric) *string {
	if !n.Valid {
		return nil
	}

	val, err := n.Value()
	if err != nil {
		panic(err)
	}

	if val == nil {
		return nil
	}

	s := val.(string)
	return &s
}

func TimestamptzToTime(t pgtype.Timestamptz) *time.Time {
	if !t.Valid {
		return nil
	}

	val := t.Time
	return &val
}
