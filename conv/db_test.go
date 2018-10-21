package conv

import (
	"database/sql"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestConvertFloat(t *testing.T) {
	inputValid := sql.NullFloat64{1, true}
	actualOutput := ConvertFloat(inputValid)
	assert.EqualValues(t, 1, actualOutput)

	inputInValid := sql.NullFloat64{1, false}
	actualOutput = ConvertFloat(inputInValid)
	assert.EqualValues(t, 0, actualOutput)
}

func TestConvertInt(t *testing.T) {
	inputValid := sql.NullInt64{1, true}
	actualOutput := ConvertInt(inputValid)
	assert.EqualValues(t, 1, actualOutput)

	inputInValid := sql.NullInt64{1, false}
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
	inputValid := sql.NullBool{true, true}
	actualOutput := ConvertBool(inputValid)
	assert.True(t, *actualOutput)

	inputInValid := sql.NullBool{true, false}
	actualOutput = ConvertBool(inputInValid)
	assert.Nil(t, actualOutput)
}

func TestConvertString(t *testing.T) {
	inputValid := sql.NullString{"one", true}
	actualOutput := ConvertString(inputValid)
	assert.EqualValues(t, "one", actualOutput)

	inputInValid := sql.NullString{"one", false}
	actualOutput = ConvertString(inputInValid)
	assert.EqualValues(t, "", actualOutput)
}