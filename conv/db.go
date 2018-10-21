package conv

import (
	"database/sql"
	"time"
)

const MYSQL_TIME_FORMAT="2006-01-02 15:04:05"

func ConvertFloat(input sql.NullFloat64) float64 {
	var floatVal float64 = 0
	if input.Valid {
		floatVal = input.Float64
	}

	return floatVal
}

func ConvertInt(input sql.NullInt64) int64 {
	var intVal int64 = 0
	if input.Valid {
		intVal = input.Int64
	}

	return intVal
}

func FormatTimePointer(input *time.Time) interface{} {
	if input == nil {
		return nil
	}

	return input.Format(MYSQL_TIME_FORMAT)
}

func ConvertBool(input sql.NullBool) (output *bool) {
	output = nil
	if input.Valid {
		output = &input.Bool
	}

	return
}

func ConvertString(input sql.NullString) (output string) {
	output = ""
	if input.Valid {
		output = input.String
	}

	return
}
