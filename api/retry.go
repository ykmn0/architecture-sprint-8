package main

import (
	"fmt"
	"time"
)

// Retry performs an operation with retry attempts and exponential backoff.
func Retry(attempts int, sleep time.Duration, function func() error) error {
	for i := 0; ; i++ {
		// Execute the function
		err := function()
		if err == nil {
			return nil
		}

		// If the maximum number of attempts is reached, return the last error
		if i >= (attempts - 1) {
			return fmt.Errorf("after %d attempts, last error: %s", attempts, err)
		}

		// Wait before the next attempt with an increasing backoff time
		time.Sleep(sleep)
		sleep = sleep * 2 // Exponential backoff
	}
}
