package time

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGetDateAtDayStart(t *testing.T) {
	inputDate, err := time.Parse("2006-01-02T15:04:05", "2001-01-01T11:00:00")
	assert.NoError(t, err)

	expectedDateBeginning, err := time.Parse("2006-01-02T15:04:05", "2001-01-01T00:00:00")
	assert.NoError(t, err)
	actualDateBeginning := GetDateAtDayStart(inputDate)
	assert.Equal(t, expectedDateBeginning.UTC(), actualDateBeginning)

	expectedDateEnd, err := time.Parse("2006-01-02T15:04:05", "2001-01-01T23:59:59.999999999")
	assert.NoError(t, err)
	actualDateEnd := GetDateAtDayEnd(inputDate)
	assert.Equal(t, expectedDateEnd.UTC(), actualDateEnd)
}
