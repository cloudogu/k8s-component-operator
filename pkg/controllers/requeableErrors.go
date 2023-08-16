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
func (gre *genericRequeueableError) GetRequeueTime(requeueTimeNanos time.Duration) time.Duration {
	return getRequeueTime(requeueTimeNanos)
}

// Unwrap returns the root error.
func (gre *genericRequeueableError) Unwrap() error {
	return gre.err
}

type dependencyUnsatisfiedError struct {
	err error
}

// Error returns the string representation of the wrapped error.
func (due *dependencyUnsatisfiedError) Error() string {
	return fmt.Sprintf("one or more dependencies are not satisfied: %s", due.err.Error())
}

// GetRequeueTime returns the time until the component should be requeued.
func (due *dependencyUnsatisfiedError) GetRequeueTime(requeueTimeNanos time.Duration) time.Duration {
	return getRequeueTime(requeueTimeNanos)
}

// Unwrap returns the root error.
func (due *dependencyUnsatisfiedError) Unwrap() error {
	return due.err
}

func getRequeueTime(currentRequeueTime time.Duration) time.Duration {
	const initialRequeueTime = 15 * time.Second
	const linearCutoffThreshold6Hours = 6 * time.Hour

	if currentRequeueTime == 0 {
		return initialRequeueTime
	}

	nextRequeueTime := currentRequeueTime * 2

	if nextRequeueTime >= linearCutoffThreshold6Hours {
		return linearCutoffThreshold6Hours
	}

	return nextRequeueTime
}
