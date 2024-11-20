package client

import (
	"bytes"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"helm.sh/helm/v3/pkg/action"
	"os"
	"testing"
	"time"
)

func Test_provider_newInstall(t *testing.T) {
	// given
	sut := &provider{
		Configuration: &action.Configuration{},
		plainHttp:     true,
		insecureTls:   true,
	}

	// when
	result := sut.newInstall()

	// then
	assert.NotEmpty(t, result.raw())
	assert.True(t, result.raw().PlainHTTP)
	assert.True(t, result.raw().InsecureSkipTLSverify)
}

func Test_provider_newInstall_providerOptionsNotSet(t *testing.T) {
	// given
	sut := &provider{Configuration: &action.Configuration{}}

	// when
	result := sut.newInstall()

	// then
	assert.NotEmpty(t, result.raw())
	assert.False(t, result.raw().PlainHTTP)
	assert.False(t, result.raw().InsecureSkipTLSverify)
}

func Test_provider_newUpgrade(t *testing.T) {
	// given
	sut := &provider{
		Configuration: &action.Configuration{},
		plainHttp:     true,
		insecureTls:   true,
	}

	// when
	result := sut.newUpgrade()

	// then
	assert.NotEmpty(t, result.raw())
	assert.True(t, result.raw().PlainHTTP)
	assert.True(t, result.raw().InsecureSkipTLSverify)
}
func Test_provider_newUpgrade_providerOptionsNotSet(t *testing.T) {
	// given
	sut := &provider{Configuration: &action.Configuration{}}

	// when
	result := sut.newUpgrade()

	// then
	assert.NotEmpty(t, result.raw())
	assert.False(t, result.raw().PlainHTTP)
	assert.False(t, result.raw().InsecureSkipTLSverify)
}

func Test_provider_newUninstall(t *testing.T) {
	// given
	sut := &provider{Configuration: &action.Configuration{}}

	// when
	result := sut.newUninstall()

	// then
	assert.NotEmpty(t, result.raw())
}

func Test_provider_newLocateChart(t *testing.T) {
	// given
	sut := &provider{Configuration: &action.Configuration{}}

	// when
	result := sut.newLocateChart()

	// then
	assert.NotEmpty(t, result)
}

func Test_provider_newListReleases(t *testing.T) {
	// given
	sut := &provider{Configuration: &action.Configuration{}}

	// when
	result := sut.newListReleases()

	// then
	assert.NotEmpty(t, result.raw())
}

func Test_provider_newGetReleaseValues(t *testing.T) {
	// given
	sut := &provider{Configuration: &action.Configuration{}}

	// when
	result := sut.newGetReleaseValues()

	// then
	assert.NotEmpty(t, result.raw())
}

func Test_provider_newGetRelease(t *testing.T) {
	// given
	sut := &provider{Configuration: &action.Configuration{}}

	// when
	result := sut.newGetRelease()

	// then
	assert.NotEmpty(t, result.raw())
}

func Test_provider_newRollbackRelease(t *testing.T) {
	// given
	sut := &provider{Configuration: &action.Configuration{}}

	// when
	result := sut.newRollbackRelease()

	// then
	assert.NotEmpty(t, result.raw())
}

func Test_readRollbackReleaseTimeoutMinsEnv(t *testing.T) {
	tests := []struct {
		name        string
		setEnvVar   bool
		envVarValue string
		want        time.Duration
		wantLogs    bool
		wantedLogs  string
		logLevel    logrus.Level
	}{
		{
			name:       "Environment variable not set",
			setEnvVar:  false,
			want:       time.Duration(15),
			wantLogs:   true,
			wantedLogs: "failed to read ROLLBACK_RELEASE_TIMEOUT_MINS environment variable, using default value of 15",
			logLevel:   logrus.DebugLevel,
		},
		{
			name:        "Environment variable not set correctly",
			setEnvVar:   true,
			envVarValue: "15//",
			want:        time.Duration(15),
			wantLogs:    true,
			wantedLogs:  "failed to parse ROLLBACK_RELEASE_TIMEOUT_MINS environment variable, using default value of 15",
			logLevel:    logrus.WarnLevel,
		},
		{
			name:        "read negative environment variable",
			setEnvVar:   true,
			envVarValue: "-20",
			want:        time.Duration(15),
			wantLogs:    true,
			wantedLogs:  "parsed value (-20) is smaller than 0, using default value of 15",
			logLevel:    logrus.WarnLevel,
		},
		{
			name:        "Successfully read environment variable",
			setEnvVar:   true,
			envVarValue: "20",
			want:        time.Duration(20),
			wantLogs:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setEnvVar {
				err := os.Setenv(rollbackReleaseTimeoutMinsEnv, tt.envVarValue)
				require.NoError(t, err)
			}
			var result = time.Duration(0)

			var logOutput bytes.Buffer

			originalOutput := logrus.StandardLogger().Out
			originalLevel := logrus.StandardLogger().Level
			if tt.wantLogs {
				logrus.StandardLogger().SetOutput(&logOutput)
				logrus.StandardLogger().SetLevel(tt.logLevel)
			}

			result = readRollbackReleaseTimeoutMinsEnv()

			logrus.StandardLogger().SetOutput(originalOutput)
			logrus.StandardLogger().SetLevel(originalLevel)

			logs := logOutput.String()

			assert.Equalf(t, tt.want, result, "readRollbackReleaseTimeoutMinsEnv()")

			if tt.wantLogs {
				assert.Contains(t, logs, tt.wantedLogs)
			}
		})
	}
}
