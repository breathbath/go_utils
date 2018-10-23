package time

import (
	"github.com/breathbath/go_utils/utils/math"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGetUnixTimestampMilliseconds(t *testing.T) {
	actualTimestamp := GetUnixTimestampMilliseconds()
	numbersCount := math.CountDigits(actualTimestamp)
	assert.EqualValues(
		t,
		13,
		numbersCount,
		"Unix timestamp %d in ms should have 13 digits rather than %d",
		actualTimestamp,
		numbersCount,
	)
}

func TestGetTimeFromTimestampMilliseconds(t *testing.T) {
	var timestampMs int64 = 978346800000
	actualTime := GetTimeFromTimestampMilliseconds(timestampMs)
	expectedTime, _ := time.Parse("2006-01-02T15:04:05", "2001-01-01T11:00:00")
	assert.Equal(t, expectedTime.UTC(), actualTime)
}