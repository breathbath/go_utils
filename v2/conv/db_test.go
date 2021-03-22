package conv

import (
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConvertFloat(t *testing.T) {
	inputValid := sql.NullFloat64{Float64: 1, Valid: true}
	actualOutput := ConvertFloat(inputValid)
	assert.EqualValues(t, 1, actualOutput)

	inputInValid := sql.NullFloat64{Float64: 1, Valid: false}
	actualOutput = ConvertFloat(inputInValid)
	assert.EqualValues(t, 0, actualOutput)
}

func TestConvertInt(t *testing.T) {
	inputValid := sql.NullInt64{Int64: 1, Valid: true}
	actualOutput := ConvertInt(inputValid)
	assert.EqualValues(t, 1, actualOutput)

	inputInValid := sql.NullInt64{Int64: 1, Valid: false}
	actualOutput = ConvertInt(inputInValid)
	assert.EqualValues(t, 0, actualOutput)
}

func TestFormatTimePointer(t *testing.T) {
	tm, err := time.Parse("2006-01-02T15:04:05", "2001-01-01T01:02:30")
	assert.NoError(t, err)

	actualTime := FormatTimePointer(&tm)
	assert.Equal(t, "2001-01-01 01:02:30", actualTime)

	nilTime := FormatTimePointer(nil)
	assert.Nil(t, nilTime)
}

func TestConvertBool(t *testing.T) {
	inputValid := sql.NullBool{Bool: true, Valid: true}
	actualOutput := ConvertBool(inputValid)
	assert.True(t, *actualOutput)

	inputInValid := sql.NullBool{Bool: true, Valid: false}
	actualOutput = ConvertBool(inputInValid)
	assert.Nil(t, actualOutput)
}

func TestConvertString(t *testing.T) {
	inputValid := sql.NullString{String: "one", Valid: true}
	actualOutput := ConvertString(inputValid)
	assert.EqualValues(t, "one", actualOutput)

	inputInValid := sql.NullString{String: "one", Valid: false}
	actualOutput = ConvertString(inputInValid)
	assert.EqualValues(t, "", actualOutput)
}
