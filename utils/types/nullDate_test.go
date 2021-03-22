package types

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNullDateToJsonConversion(t *testing.T) {
	expectedTime, err := time.Parse("2006-01-02", "2001-01-01")
	assert.NoError(t, err)

	nullDate := NullDate{}

	jsonDate, err := nullDate.MarshalJSON()
	assert.NoError(t, err)
	assert.Equal(t, "null", string(jsonDate))

	nullDate.Time = expectedTime
	nullDate.Valid = true

	jsonDate, err = nullDate.MarshalJSON()
	assert.NoError(t, err)
	assert.Equal(t, `"2001-01-01"`, string(jsonDate))
}

func TestNullDateFromJsonConversion(t *testing.T) {
	jsonTime := `"2001-01-01"`
	nullDate := NullDate{}

	err := nullDate.UnmarshalJSON([]byte(jsonTime))
	assert.NoError(t, err)

	expectedTime, err := time.Parse("2006-01-02", "2001-01-01")
	assert.NoError(t, err)

	assert.True(t, nullDate.Valid)
	assert.Equal(t, expectedTime, nullDate.Time)

	err = nullDate.UnmarshalJSON([]byte("invalid date"))
	assert.EqualError(t, err, "parsing time \"invalid date\" as \"\"2006-01-02\"\": cannot parse \"invalid date\" as \"\"\"")

	err = nullDate.UnmarshalJSON([]byte(`"2001-01-01 11:00:00"`))
	assert.EqualError(t, err, "parsing time \"\"2001-01-01 11:00:00\"\" as \"\"2006-01-02\"\": cannot parse \" 11:00:00\"\" as \"\"\"")

	err = nullDate.UnmarshalJSON([]byte(`""`))
	assert.NoError(t, err)
	assert.False(t, nullDate.Valid)

	err = nullDate.UnmarshalJSON([]byte("null"))
	assert.NoError(t, err)
	assert.False(t, nullDate.Valid)
}
