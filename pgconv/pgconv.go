package pgconv

import (
	"log"
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

func TimeToTimestamptz(t time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{
		Time:             t,
		InfinityModifier: pgtype.Finite,
		Valid:            true,
	}
}

func StringToNumeric(s string) pgtype.Numeric {
	var numeric pgtype.Numeric

	if err := numeric.Scan(s); err != nil {
		log.Panic(err)
	}

	return numeric
}

func Int64ToInt8(i int64) pgtype.Int8 {
	return pgtype.Int8{Int64: i, Valid: true}
}
