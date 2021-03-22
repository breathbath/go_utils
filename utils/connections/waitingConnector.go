// Package connections allows to establish connections to remote resources which are not immediately available
// by repeating attempts with sleeping intervals
package connections

import (
	"errors"
	"fmt"
	"log"
	"time"
)

type Sleeper func(sleepTime time.Duration)

var currentSleeper = func(sleepTime time.Duration) {
	time.Sleep(sleepTime)
}

// WaitingConnectorIterativeDelayDuration global var allows to tweak sleeping interval between failed connections
var WaitingConnectorIterativeDelayDuration = time.Second

// SetSleeper global func allows to replace sleeping logic with non-sleeping implementation for tests reasons
func SetSleeper(sleeper Sleeper) {
	currentSleeper = sleeper
}

// WaitForConnection tries to connect to a remote resource with resourceName using conn function resourceCallback
// stopping after maxConnAttempts. outputFunc is used to output the info about attempts to std out which can be
// replaced with some custom output func
func WaitForConnection(
	maxConnAttempts int,
	resourceName string,
	resourceCallback func() (interface{}, error),
	outputFunc func(msg string, err error),
) (interface{}, error) {
	if outputFunc == nil {
		outputFunc = func(msg string, err error) {
			if msg != "" {
				log.Println(msg)
			}

			if err != nil {
				log.Println(err)
			}
		}
	}

	for i := 0; i < maxConnAttempts; i++ {
		res, err := resourceCallback()

		if err == nil {
			return res, nil
		}

		errorToSend := fmt.Errorf("%s connection error: %s", resourceName, err.Error())

		outputFunc("", errorToSend)

		sleepTime := WaitingConnectorIterativeDelayDuration * time.Duration(i+1)
		outputFunc(fmt.Sprintf("Trying to reconnect to %s in %.0f s\n", resourceName, sleepTime.Seconds()), nil)
		currentSleeper(sleepTime)
	}

	return nil, errors.New("Was not able to connect to " + resourceName)
}
