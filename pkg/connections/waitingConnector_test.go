package connections

import (
	"errors"
	"testing"
	"time"

	testing2 "github.com/breathbath/go_utils/v3/pkg/testing"
	"github.com/stretchr/testify/assert"
)

var mockedSleeper = func(timeToSleep time.Duration) {}

func TestSuccessConnectionAfterFirstAttempt(t *testing.T) {
	SetSleeper(mockedSleeper)
	connector := func() (interface{}, error) {
		return 1, nil
	}

	outputFunc := func(msg string, err error) {}

	res, err := WaitForConnection(10, "", connector, outputFunc)

	assert.NoError(t, err)
	assert.Equal(t, 1, res)
}

func TestSleepingIntervals(t *testing.T) {
	connector := func() (interface{}, error) {
		return nil, errors.New("some error")
	}
	outputFunc := func(msg string, err error) {}

	sleepingIntervals := []float64{}
	sleepFunc := func(timeToSleep time.Duration) {
		sleepingIntervals = append(sleepingIntervals, timeToSleep.Seconds())
	}
	SetSleeper(sleepFunc)

	_, err := WaitForConnection(10, "", connector, outputFunc)
	assert.Error(t, err)

	expectedSleepingIntervals := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	assert.Equal(t, expectedSleepingIntervals, sleepingIntervals)
}

func TestSuccessConnectionAfterXAttempts(t *testing.T) {
	SetSleeper(mockedSleeper)
	callsCount := 0
	connector := func() (interface{}, error) {
		callsCount++
		if callsCount > 1 {
			return 1, nil
		}

		return nil, errors.New("connection error")
	}

	allOutputs := []string{}

	outputFunc := func(msg string, err error) {
		if msg != "" {
			allOutputs = append(allOutputs, msg)
		}
		if err != nil {
			allOutputs = append(allOutputs, err.Error())
		}
	}

	res, err := WaitForConnection(2, "SomeConn", connector, outputFunc)
	assert.NoError(t, err)
	assert.Equal(t, 1, res)
	assert.Equal(t, 2, callsCount)

	expectedOutputs := []string{
		"SomeConn connection error: connection error",
		"Trying to reconnect to SomeConn in 1 s\n",
	}
	assert.Equal(t, expectedOutputs, allOutputs)
}

func TestConnectionFailure(t *testing.T) {
	SetSleeper(mockedSleeper)
	connector := func() (interface{}, error) {
		return nil, errors.New("connection error")
	}
	outputFunc := func(msg string, err error) {}

	res, err := WaitForConnection(2, "SomeConn", connector, outputFunc)

	assert.EqualError(t, err, "Was not able to connect to SomeConn")
	assert.Nil(t, res)
}

func TestNoOutputFuncProvided(t *testing.T) {
	SetSleeper(mockedSleeper)
	connector := func() (interface{}, error) {
		return 1, errors.New("some error")
	}

	output := testing2.CaptureOutput(func() {
		_, err := WaitForConnection(2, "RemoteHost", connector, nil)
		assert.Error(t, err)
	})

	assert.Contains(t, output, "some error")
	assert.Contains(t, output, "Trying to reconnect to RemoteHost in 1 s")
	assert.Contains(t, output, "Trying to reconnect to RemoteHost in 2 s")
}
