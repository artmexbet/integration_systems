package utills

import (
	"github.com/jackc/pgx/v5/pgtype"
	"strconv"
)

func GetPtr[T any](v T) *T {
	return &v
}

func GetVal[T any](v *T) T {
	if v == nil {
		var zero T
		return zero
	}
	return *v
}

func ParseStringToInt(s string) int {
	result, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return result
}

func Int32ToPgInt4(v int32) pgtype.Int4 {
	res := pgtype.Int4{}
	err := res.Scan(v)
	if err != nil {
		return pgtype.Int4{}
	}
	return res
}

func StringToPgText(s string) pgtype.Text {
	res := pgtype.Text{}
	err := res.Scan(s)
	if err != nil {
		return pgtype.Text{}
	}
	return res
}
