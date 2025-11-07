package controllers

import (
	"errors"
	"testing"
	"time"

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

func Test_genericRequeueableError_GetRequeueTime(t *testing.T) {
	type args struct {
		requeueTime        time.Duration
		defaultRequeueTime time.Duration
	}
	tests := []struct {
		name string
		args args
		want time.Duration
	}{
		{"always return defaultRequeueTime case 1", args{0 * time.Second, 5 * time.Second}, 5 * time.Second},
		{"always return defaultRequeueTime case 2", args{15 * time.Second, 5 * time.Second}, 5 * time.Second},
		{"always return defaultRequeueTime case 3", args{30 * time.Second, 5 * time.Second}, 5 * time.Second},
		{"always return defaultRequeueTime case 3", args{256 * time.Second, 5 * time.Second}, 5 * time.Second},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			due := &genericRequeueableError{}
			assert.Equalf(t, tt.want, due.GetRequeueTime(tt.args.requeueTime, tt.args.defaultRequeueTime), "getRequeueTime(%v)", tt.args.requeueTime)
		})
	}
}
