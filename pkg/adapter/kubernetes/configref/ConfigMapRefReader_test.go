package configref

import (
	"context"
	"testing"

	v1 "github.com/cloudogu/k8s-component-lib/api/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

var testCtx = context.Background()

const componentName = "testComponent"

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

func TestConfigMapRefReader_GetSystemValues(t *testing.T) {
	t.Run("component is nil", func(t *testing.T) {
		configMapMock := newMockConfigMapClient(t)
		refReader := NewConfigMapRefReader(configMapMock)

		result, err := refReader.GetSystemValues(testCtx, nil)
		require.NoError(t, err)
		assert.Equal(t, "", result)
	})
	t.Run("no configmap found", func(t *testing.T) {
		component := &v1.Component{ObjectMeta: metav1.ObjectMeta{Name: componentName}}
		configMapMock := newMockConfigMapClient(t)
		cmLabelSelector := labels.ValidatedSetSelector{
			"k8s.cloudogu.com/component.config": component.Name,
		}.String()
		configMapMock.EXPECT().List(testCtx, metav1.ListOptions{LabelSelector: cmLabelSelector}).Return(&corev1.ConfigMapList{}, nil)
		refReader := NewConfigMapRefReader(configMapMock)

		result, err := refReader.GetSystemValues(testCtx, component)
		require.NoError(t, err)
		assert.Equal(t, "", result)
	})
	t.Run("configmap contains no keys", func(t *testing.T) {
		cm := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name: "testComponent-config",
				Labels: map[string]string{
					"k8s.cloudogu.com/component.config": componentName,
				},
			},
		}
		component := &v1.Component{ObjectMeta: metav1.ObjectMeta{Name: componentName}}
		configMapMock := newMockConfigMapClient(t)
		cmLabelSelector := labels.ValidatedSetSelector{
			"k8s.cloudogu.com/component.config": component.Name,
		}.String()
		configMapMock.EXPECT().List(testCtx, metav1.ListOptions{LabelSelector: cmLabelSelector}).Return(&corev1.ConfigMapList{Items: []corev1.ConfigMap{*cm}}, nil)
		refReader := NewConfigMapRefReader(configMapMock)

		result, err := refReader.GetSystemValues(testCtx, component)
		require.Error(t, err)
		assert.ErrorContains(t, err, "key `values` does not exist in configmap")
		assert.Equal(t, "", result)
	})
	t.Run("successfully get configmap values", func(t *testing.T) {
		cm := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name: "testComponent-config",
				Labels: map[string]string{
					"k8s.cloudogu.com/component.config": componentName,
				},
			},
			Data: map[string]string{
				"values": "test",
			},
		}
		component := &v1.Component{ObjectMeta: metav1.ObjectMeta{Name: componentName}}
		configMapMock := newMockConfigMapClient(t)
		cmLabelSelector := labels.ValidatedSetSelector{
			"k8s.cloudogu.com/component.config": component.Name,
		}.String()
		configMapMock.EXPECT().List(testCtx, metav1.ListOptions{LabelSelector: cmLabelSelector}).Return(&corev1.ConfigMapList{Items: []corev1.ConfigMap{*cm}}, nil)
		refReader := NewConfigMapRefReader(configMapMock)

		result, err := refReader.GetSystemValues(testCtx, component)
		require.NoError(t, err)
		assert.Equal(t, "test", result)
	})
}
