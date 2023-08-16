package controllers

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_dependencyUnsatisfiedError_GetRequeueTime(t *testing.T) {
	type args struct {
		requeueTime time.Duration
	}
	tests := []struct {
		name string
		args args
		want time.Duration
	}{
		// double the value until the threshold jumps in
		{"1st interval", args{0 * time.Second}, 15 * time.Second},
		{"2nd interval", args{15 * time.Second}, 30 * time.Second},
		{"3rd interval", args{30 * time.Second}, 1 * time.Minute},
		{"11th interval", args{128 * time.Minute}, 256 * time.Minute},
		{"cutoff interval ", args{256 * time.Minute}, 6 * time.Hour},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			due := &dependencyUnsatisfiedError{}
			assert.Equalf(t, tt.want, due.GetRequeueTime(tt.args.requeueTime), "getRequeueTime(%v)", tt.args.requeueTime)
		})
	}
}

func Test_dependencyUnsatisfiedError_Unwrap(t *testing.T) {
	testErr1 := assert.AnError
	testErr2 := errors.New("test")
	inputErr := errors.Join(testErr1, testErr2)

	sut := &dependencyUnsatisfiedError{inputErr}

	// when
	actualErr := sut.Unwrap()

	// then
	require.Error(t, sut)
	require.Error(t, actualErr)
	assert.ErrorIs(t, actualErr, testErr1)
	assert.ErrorIs(t, actualErr, testErr2)
}

func Test_dependencyUnsatisfiedError_Error(t *testing.T) {
	sut := &dependencyUnsatisfiedError{assert.AnError}
	expected := "one or more dependencies are not satisfied: assert.AnError general error for testing"
	assert.Equal(t, expected, sut.Error())
}

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
