package connections

import (
	"fmt"
	"time"

	"github.com/breathbath/go_utils/v3/pkg/errs"
)

// This example tries to connect to a non immediately available resource outputting attempts info
func ExampleWaitForConnection() {
	attemptsToConnect := 2
	connFunc := func() (interface{}, error) {
		if attemptsToConnect > 0 {
			attemptsToConnect--
			return nil, fmt.Errorf("cannot connect to api.com")
		}
		return "SomeConnectionObject", nil
	}

	SetSleeper(func(sleepTime time.Duration) {
		// we don't sleep at all to not slow down the test
	})

	res, err := WaitForConnection(
		3,
		"api.com",
		connFunc,
		func(msg string, err error) {
			if err != nil {
				fmt.Printf("Error:%v", err)
			}
			fmt.Printf("%s", msg)
		},
	)
	errs.FailOnError(err)

	fmt.Printf("This is my resource after 2 connection attempts: '%s'", res.(string))

	// Output:
	// Error:api.com connection error: cannot connect to api.comTrying to reconnect to api.com in 1 s
	// Error:api.com connection error: cannot connect to api.comTrying to reconnect to api.com in 2 s
	// This is my resource after 2 connection attempts: 'SomeConnectionObject'
}
