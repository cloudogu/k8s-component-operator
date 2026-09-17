package controllers

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_genericRequeueableError_Unwrap(t *testing.T) {
	testErr1 := assert.AnError
	testErr2 := errors.New("test")
	inputErr := errors.Join(testErr1, testErr2)

	sut := &genericRequeueableError{"oh noez", inputErr}

	// when
	actualErr := sut.Unwrap()

	// then
	require.Error(t, sut)
	require.Error(t, actualErr)
	assert.ErrorIs(t, actualErr, testErr1)
	assert.ErrorIs(t, actualErr, testErr2)
}

func Test_genericRequeueableError_Error(t *testing.T) {
	sut := &genericRequeueableError{"oh noez", assert.AnError}
	expected := "oh noez: " + assert.AnError.Error()
	assert.Equal(t, expected, sut.Error())
}
