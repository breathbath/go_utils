package types

import (
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestJsonConversion(t *testing.T) {
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
