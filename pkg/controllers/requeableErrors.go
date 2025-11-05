package controllers

import (
	"fmt"
	"time"
)

type genericRequeueableError struct {
	errMsg string
	err    error
}

// Error returns the string representation of the wrapped error.
func (gre *genericRequeueableError) Error() string {
	return fmt.Sprintf("%s: %s", gre.errMsg, gre.err.Error())
}

// GetRequeueTime returns the time until the component should be requeued.
func (gre *genericRequeueableError) GetRequeueTime(requeueTimeNanos time.Duration, defaultRequeueTimeNanos time.Duration) time.Duration {
	return getRequeueTime(requeueTimeNanos, defaultRequeueTimeNanos)
}

// Unwrap returns the root error.
func (gre *genericRequeueableError) Unwrap() error {
	return gre.err
}

func getRequeueTime(currentRequeueTime time.Duration, defaultRequeueTimeNanos time.Duration) time.Duration {

	if currentRequeueTime == 0 {
		return defaultRequeueTimeNanos
	}

	nextRequeueTime := defaultRequeueTimeNanos

	return nextRequeueTime
}
