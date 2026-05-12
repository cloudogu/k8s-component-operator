package configref

import (
	"context"
	"testing"

	v1 "github.com/cloudogu/k8s-component-lib/api/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var testCtx = context.Background()

var (
	valuesConfigMap = &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name: "valuesConfigMap",
		},
		Data: map[string]string{"values": "user1"},
	}
	valuesConfigMapWithMissingKey = &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name: "valuesConfigMap",
		},
		Data: map[string]string{},
	}
)

func TestConfigMapRefReader_GetValues(t *testing.T) {
	t.Run("nothing to load", func(t *testing.T) {
		configMapMock := newMockConfigMapClient(t)
		refReader := NewConfigMapRefReader(configMapMock)

		result, err := refReader.GetValues(testCtx, &v1.Reference{})
		require.NoError(t, err)
		assert.Equal(t, "", result)
	})
	t.Run("nothing to load with nil input", func(t *testing.T) {
		configMapMock := newMockConfigMapClient(t)
		refReader := NewConfigMapRefReader(configMapMock)

		result, err := refReader.GetValues(testCtx, nil)
		require.NoError(t, err)
		assert.Equal(t, "", result)
	})
	t.Run("load config map with key", func(t *testing.T) {
		configMapMock := newMockConfigMapClient(t)
		configMapMock.EXPECT().
			Get(testCtx, "valuesConfigMap", metav1.GetOptions{}).
			Return(valuesConfigMap, nil)

		refReader := NewConfigMapRefReader(configMapMock)

		result, err := refReader.GetValues(testCtx, &v1.Reference{
			Name: "valuesConfigMap",
			Key:  "values",
		})
		require.NoError(t, err)
		assert.Equal(t, "user1", result)
	})
	t.Run("try load missing config map", func(t *testing.T) {
		configMapMock := newMockConfigMapClient(t)
		configMapMock.EXPECT().
			Get(testCtx, "valuesConfigMap", metav1.GetOptions{}).
			Return(nil, assert.AnError)

		refReader := NewConfigMapRefReader(configMapMock)

		result, err := refReader.GetValues(testCtx, &v1.Reference{
			Name: "valuesConfigMap",
			Key:  "values",
		})
		require.Error(t, err)
		assert.Equal(t, "", result)
	})
	t.Run("try load missing config map key", func(t *testing.T) {
		configMapMock := newMockConfigMapClient(t)
		configMapMock.EXPECT().
			Get(testCtx, "valuesConfigMap", metav1.GetOptions{}).
			Return(valuesConfigMapWithMissingKey, nil)

		refReader := NewConfigMapRefReader(configMapMock)

		result, err := refReader.GetValues(testCtx, &v1.Reference{
			Name: "valuesConfigMap",
			Key:  "values",
		})
		require.Error(t, err)
		assert.Equal(t, "", result)
	})
}
