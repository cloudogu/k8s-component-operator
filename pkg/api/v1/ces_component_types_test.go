package v1

import (
	"bytes"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestComponent_GetHelmChartSpec(t *testing.T) {
	type fields struct {
		TypeMeta   v1.TypeMeta
		ObjectMeta v1.ObjectMeta
		Spec       ComponentSpec
		Status     ComponentStatus
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{name: "should use deployNamespace if specified", fields: fields{Spec: ComponentSpec{DeployNamespace: "longhorn"}}, want: "longhorn"},
		{name: "should use regular namespace if no deployNamespace if specified", fields: fields{ObjectMeta: v1.ObjectMeta{Namespace: "ecosystem"}, Spec: ComponentSpec{DeployNamespace: ""}}, want: "ecosystem"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Component{
				TypeMeta:   tt.fields.TypeMeta,
				ObjectMeta: tt.fields.ObjectMeta,
				Spec:       tt.fields.Spec,
				Status:     tt.fields.Status,
			}
			if got := c.GetHelmChartSpec().Namespace; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetHelmChartSpec() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_readHelmClientTimeoutMinsEnv(t *testing.T) {
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
			wantedLogs: "failed to read HELM_CLIENT_TIMEOUT_MINS environment variable, using default value of 15",
			logLevel:   logrus.DebugLevel,
		},
		{
			name:        "Environment variable not set correctly",
			setEnvVar:   true,
			envVarValue: "15//",
			want:        time.Duration(15),
			wantLogs:    true,
			wantedLogs:  "failed to parse HELM_CLIENT_TIMEOUT_MINS environment variable, using default value of 15",
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
				err := os.Setenv(helmClientTimeoutMinsEnv, tt.envVarValue)
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

			result = readHelmClientTimeoutMinsEnv()

			logrus.StandardLogger().SetOutput(originalOutput)
			logrus.StandardLogger().SetLevel(originalLevel)

			logs := logOutput.String()

			assert.Equalf(t, tt.want, result, "readHelmClientTimeoutMinsEnv()")

			if tt.wantLogs {
				assert.Contains(t, logs, tt.wantedLogs)
			}
		})
	}
}
