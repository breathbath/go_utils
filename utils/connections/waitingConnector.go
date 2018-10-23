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

func SetSleeper(sleeper Sleeper) {
	currentSleeper = sleeper
}

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

	for i := 0; i < int(maxConnAttempts); i++ {
		res, err := resourceCallback()

		if err == nil {
			return res, nil
		}

		errorToSend := fmt.Errorf("%s connection error: %s", resourceName, err.Error())

		outputFunc("", errorToSend)

		sleepTime := time.Second * time.Duration(i+1)
		outputFunc(fmt.Sprintf("Trying to reconnect to %s in %.0f s\n", resourceName, sleepTime.Seconds()), nil)
		currentSleeper(sleepTime)
	}

	return nil, errors.New("Was not able to connect to " + resourceName)
}
