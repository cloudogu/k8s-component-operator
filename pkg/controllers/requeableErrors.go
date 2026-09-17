package controllers

import (
	"fmt"
)

type genericRequeueableError struct {
	errMsg string
	err    error
}

// Error returns the string representation of the wrapped error.
func (gre *genericRequeueableError) Error() string {
	return fmt.Sprintf("%s: %s", gre.errMsg, gre.err.Error())
}

// Unwrap returns the root error.
func (gre *genericRequeueableError) Unwrap() error {
	return gre.err
}
