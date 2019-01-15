package types

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNullDateFromJsonConversion(t *testing.T) {
	jsonTime := `"2001-01-01"`
	nullTime := NullDate{}

	err := nullTime.UnmarshalJSON([]byte(jsonTime))
	assert.NoError(t, err)

	expectedTime, err := time.Parse("2006-01-02", "2001-01-01")
	assert.NoError(t, err)

	assert.True(t, nullTime.Valid)
	assert.Equal(t, expectedTime, nullTime.Time)

	err = nullTime.UnmarshalJSON([]byte("invalid date"))
	assert.EqualError(t, err, "parsing time \"invalid date\" as \"\"2006-01-02\"\": cannot parse \"invalid date\" as \"\"\"")

	err = nullTime.UnmarshalJSON([]byte(`"2001-01-01 11:00:00"`))
	assert.EqualError(t, err, "parsing time \"\"2001-01-01 11:00:00\"\" as \"\"2006-01-02\"\": cannot parse \" 11:00:00\"\" as \"\"\"")
}
