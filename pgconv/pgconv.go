package pgconv

import (
	"math/big"
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

func IntToNumeric(i int) pgtype.Numeric {
	bigint := big.NewInt(int64(i))

	return pgtype.Numeric{Int: bigint, Valid: true}
}

func Int64ToInt8(i int64) pgtype.Int8 {
	return pgtype.Int8{Int64: i, Valid: true}
}
