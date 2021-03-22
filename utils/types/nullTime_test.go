package types

import (
	"testing"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

func TestToJsonConversion(t *testing.T) {
	expectedTime, err := time.Parse("2006-01-02T15:04:05", "2001-01-01T11:00:00")
	assert.NoError(t, err)

	nullTime := NullTime{
		NullTime: mysql.NullTime{Time: expectedTime, Valid: false},
	}

	jsonTime, err := nullTime.MarshalJSON()
	assert.NoError(t, err)
	assert.Equal(t, "null", string(jsonTime))

	nullTime = NullTime{
		NullTime: mysql.NullTime{Time: expectedTime, Valid: true},
	}
	jsonTime, err = nullTime.MarshalJSON()
	assert.NoError(t, err)
	assert.Equal(t, `"2001-01-01T11:00:00Z"`, string(jsonTime))
}

func TestFromJsonConversion(t *testing.T) {
	jsonTime := `"2001-01-01T11:00:00Z"`
	nullTime := NullTime{}

	err := nullTime.UnmarshalJSON([]byte(jsonTime))
	assert.NoError(t, err)

	expectedTime, err := time.Parse("2006-01-02T15:04:05", "2001-01-01T11:00:00")
	assert.NoError(t, err)

	assert.True(t, nullTime.Valid)
	assert.Equal(t, expectedTime, nullTime.Time)

	err = nullTime.UnmarshalJSON([]byte("invalid date"))
	assert.EqualError(t, err, "parsing time \"invalid date\" as \"\"2006-01-02T15:04:05Z07:00\"\": cannot parse \"invalid date\" as \"\"\"")

	err = nullTime.UnmarshalJSON([]byte(`"2001-01-01 11:00:00"`))
	assert.EqualError(t, err, "parsing time \"\"2001-01-01 11:00:00\"\" as \"\"2006-01-02T15:04:05Z07:00\"\": cannot parse \" 11:00:00\"\" as \"T\"")

	err = nullTime.UnmarshalJSON([]byte(`""`))
	assert.NoError(t, err)
	assert.False(t, nullTime.Valid)

	err = nullTime.UnmarshalJSON([]byte("null"))
	assert.NoError(t, err)
	assert.False(t, nullTime.Valid)
}
