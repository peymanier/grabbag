package pgconv

import (
	"fmt"
	"log"
	"math/big"
	"strconv"
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

func NumericToFloat64(n pgtype.Numeric) *float64 {
	if !n.Valid {
		return nil
	}

	float8, err := n.Float64Value()
	if err != nil {
		panic(err)
	}

	if !float8.Valid {
		panic(fmt.Errorf("invalid numeric"))
	}

	return &float8.Float64
}

func Float64ToNumeric(f float64) pgtype.Numeric {
	float64Str := strconv.FormatFloat(f, 'f', -1, 64)
	return StringToNumeric(float64Str)
}

func StringToNumeric(s string) pgtype.Numeric {
	var numeric pgtype.Numeric

	if err := numeric.Scan(s); err != nil {
		log.Panic(err)
	}

	return numeric
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

func Int64ToNumeric(i int64) pgtype.Numeric {
	return pgtype.Numeric{
		Int:              big.NewInt(i),
		Exp:              0,
		NaN:              false,
		InfinityModifier: 0,
		Valid:            true,
	}
}

func Int64ToInt8(i int64) pgtype.Int8 {
	return pgtype.Int8{Int64: i, Valid: true}
}
