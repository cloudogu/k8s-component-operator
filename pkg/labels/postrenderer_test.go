package labels

import (
	"bytes"
	_ "embed"
	"fmt"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"

	yamlutil "github.com/cloudogu/k8s-component-operator/pkg/yaml"
)

//go:embed testdata/job.yaml
var jobBytes []byte

//go:embed testdata/jobWithLabels.yaml
var jobWithLabelsBytes []byte

//go:embed testdata/deployment.yaml
var deploymentBytes []byte

//go:embed testdata/deploymentWithLabels.yaml
var deploymentWithLabelsString string

//go:embed testdata/statefulSet.yaml
var statefulSetBytes []byte

//go:embed testdata/statefulSetWithLabels.yaml
var statefulSetWithLabelsString string

//go:embed testdata/daemonSet.yaml
var daemonSetBytes []byte

//go:embed testdata/daemonSetWithLabels.yaml
var daemonSetWithLabelsString string

//go:embed testdata/cronJob.yaml
var cronJobBytes []byte

//go:embed testdata/cronJobWithLabels.yaml
var cronJobWithLabelsString string

//go:embed testdata/doguOp.yaml
var doguOpBytes []byte

//go:embed testdata/doguOpWithLabels.yaml
var doguOpWithLabelsStr string

//go:embed testdata/longhorn.yaml
var longhornBytes []byte

//go:embed testdata/longhornWithLabels.yaml
var longhornWithLabelsStr string

func TestPostRenderer_Run(t *testing.T) {
	testJob := &batchv1.Job{
		TypeMeta:   metav1.TypeMeta{Kind: "Job", APIVersion: "batch/v1"},
		ObjectMeta: metav1.ObjectMeta{Name: "hello"},
		Spec: batchv1.JobSpec{Template: corev1.PodTemplateSpec{Spec: corev1.PodSpec{Containers: []corev1.Container{
			{Name: "hello", Image: "busyboxy", Command: []string{"sh", "-c", `echo "Hello, Kubernetes!" && sleep 3600`}},
		}}}},
	}

	testJobWithLabels := &batchv1.Job{
		TypeMeta: metav1.TypeMeta{Kind: "Job", APIVersion: "batch/v1"},
		ObjectMeta: metav1.ObjectMeta{
			Name: "hello",
		},
		Spec: batchv1.JobSpec{Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: map[string]string{
					"k8s.cloudogu.com/component.name":    "k8s-blueprint-operator",
					"k8s.cloudogu.com/component.version": "1.2.3-4",
				},
			},
			Spec: corev1.PodSpec{Containers: []corev1.Container{
				{Name: "hello", Image: "busyboxy", Command: []string{"sh", "-c", `echo "Hello, Kubernetes!" && sleep 3600`}},
			}},
		},
		},
	}

	type fields struct {
		documentSplitterFn       func(t *testing.T) documentSplitter
		unstructuredSerializerFn func(t *testing.T) unstructuredSerializer
		unstructuredConverterFn  func(t *testing.T) unstructuredConverter
		serializerFn             func(t *testing.T) genericYamlSerializer
		labels                   map[string]string
	}
	tests := []struct {
		name                     string
		fields                   fields
		renderedManifests        *bytes.Buffer
		wantModifiedManifestsStr string
		wantErr                  assert.ErrorAssertionFunc
	}{
		{
			name: "should fail to split document",
			fields: fields{
				documentSplitterFn: func(t *testing.T) documentSplitter {
					t.Helper()
					splitter := newMockDocumentSplitter(t)
					splitter.EXPECT().WithReader(mock.Anything).Return(splitter)
					splitter.EXPECT().Next().Return(false)
					splitter.EXPECT().Err().Return(assert.AnError)
					return splitter
				},
				unstructuredSerializerFn: func(t *testing.T) unstructuredSerializer {
					return newMockUnstructuredSerializer(t)
				},
				unstructuredConverterFn: func(t *testing.T) unstructuredConverter {
					return newMockUnstructuredConverter(t)
				},
				serializerFn: func(t *testing.T) genericYamlSerializer {
					return newMockGenericYamlSerializer(t)
				},
				labels: map[string]string{
					"k8s.cloudogu.com/component.name":    "k8s-blueprint-operator",
					"k8s.cloudogu.com/component.version": "1.2.3-4",
				},
			},
			renderedManifests:        bytes.NewBuffer(jobBytes),
			wantModifiedManifestsStr: "<nil>",
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, assert.AnError, i) &&
					assert.ErrorContains(t, err, "failed to split yaml document", i)
			},
		},
		{
			name: "should fail to decode object",
			fields: fields{
				documentSplitterFn: func(t *testing.T) documentSplitter {
					t.Helper()
					splitter := newMockDocumentSplitter(t)
					splitter.EXPECT().WithReader(mock.Anything).Return(splitter)
					nextCall := splitter.EXPECT().Next().Return(true).Once()
					splitter.EXPECT().Next().Return(false).NotBefore(nextCall)
					splitter.EXPECT().Bytes().Return(jobBytes)
					return splitter
				},
				unstructuredSerializerFn: func(t *testing.T) unstructuredSerializer {
					us := newMockUnstructuredSerializer(t)
					us.EXPECT().Decode(jobBytes, (*schema.GroupVersionKind)(nil), nil).Return(nil, nil, assert.AnError)
					return us
				},
				unstructuredConverterFn: func(t *testing.T) unstructuredConverter {
					return newMockUnstructuredConverter(t)
				},
				serializerFn: func(t *testing.T) genericYamlSerializer {
					return newMockGenericYamlSerializer(t)
				},
				labels: map[string]string{
					"k8s.cloudogu.com/component.name":    "k8s-blueprint-operator",
					"k8s.cloudogu.com/component.version": "1.2.3-4",
				},
			},
			renderedManifests:        bytes.NewBuffer(jobBytes),
			wantModifiedManifestsStr: "<nil>",
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, assert.AnError, i) &&
					assert.ErrorContains(t, err, "failed to parse yaml resources", i)
			},
		},
		{
			name: "should fail to convert to unstructured map",
			fields: fields{
				documentSplitterFn: func(t *testing.T) documentSplitter {
					t.Helper()
					splitter := newMockDocumentSplitter(t)
					splitter.EXPECT().WithReader(mock.Anything).Return(splitter)
					nextCall := splitter.EXPECT().Next().Return(true).Once()
					splitter.EXPECT().Next().Return(false).NotBefore(nextCall)
					splitter.EXPECT().Bytes().Return(jobBytes)
					return splitter
				},
				unstructuredSerializerFn: func(t *testing.T) unstructuredSerializer {
					us := newMockUnstructuredSerializer(t)
					us.EXPECT().Decode(jobBytes, (*schema.GroupVersionKind)(nil), nil).Return(testJob, nil, nil)
					return us
				},
				unstructuredConverterFn: func(t *testing.T) unstructuredConverter {
					uc := newMockUnstructuredConverter(t)
					uc.EXPECT().ToUnstructured(testJob).Return(nil, assert.AnError)
					return uc
				},
				serializerFn: func(t *testing.T) genericYamlSerializer {
					return newMockGenericYamlSerializer(t)
				},
				labels: map[string]string{
					"k8s.cloudogu.com/component.name":    "k8s-blueprint-operator",
					"k8s.cloudogu.com/component.version": "1.2.3-4",
				},
			},
			renderedManifests:        bytes.NewBuffer(jobBytes),
			wantModifiedManifestsStr: "<nil>",
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, assert.AnError, i) &&
					assert.ErrorContains(t, err, "failed to convert resource to unstructured object", i)
			},
		},
		{
			name: "should fail to serialize resource back to yaml",
			fields: fields{
				documentSplitterFn: func(t *testing.T) documentSplitter {
					t.Helper()
					splitter := newMockDocumentSplitter(t)
					splitter.EXPECT().WithReader(mock.Anything).Return(splitter)
					nextCall := splitter.EXPECT().Next().Return(true).Once()
					splitter.EXPECT().Next().Return(false).NotBefore(nextCall)
					splitter.EXPECT().Bytes().Return(jobBytes)
					return splitter
				},
				unstructuredSerializerFn: func(t *testing.T) unstructuredSerializer {
					us := newMockUnstructuredSerializer(t)
					us.EXPECT().Decode(jobBytes, (*schema.GroupVersionKind)(nil), nil).Return(testJob, nil, nil)
					return us
				},
				unstructuredConverterFn: func(t *testing.T) unstructuredConverter {
					uc := newMockUnstructuredConverter(t)
					uc.EXPECT().ToUnstructured(testJob).Return(jobMap(), nil).Once()
					uc.EXPECT().ToUnstructured(testJobWithLabels).Return(jobMap(), nil).Once()
					uc.EXPECT().FromUnstructured(jobMap(), &batchv1.Job{}).Run(func(um map[string]interface{}, obj interface{}) {
						job := obj.(*batchv1.Job)
						*job = *testJob
					}).Return(nil)
					return uc
				},
				serializerFn: func(t *testing.T) genericYamlSerializer {
					ys := newMockGenericYamlSerializer(t)
					ys.EXPECT().Marshal(&unstructured.Unstructured{Object: jobMapWitLabels()}).Return(nil, assert.AnError)
					return ys
				},
				labels: map[string]string{
					"k8s.cloudogu.com/component.name":    "k8s-blueprint-operator",
					"k8s.cloudogu.com/component.version": "1.2.3-4",
				},
			},
			renderedManifests:        bytes.NewBuffer(jobBytes),
			wantModifiedManifestsStr: "<nil>",
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, assert.AnError, i) &&
					assert.ErrorContains(t, err, "failed to serialize resources back to yaml", i)
			},
		},
		{
			name: "should succeed",
			fields: fields{
				documentSplitterFn: func(t *testing.T) documentSplitter {
					t.Helper()
					splitter := newMockDocumentSplitter(t)
					splitter.EXPECT().WithReader(mock.Anything).Return(splitter)
					nextCall := splitter.EXPECT().Next().Return(true).Once()
					splitter.EXPECT().Next().Return(false).NotBefore(nextCall)
					splitter.EXPECT().Bytes().Return(jobBytes)
					splitter.EXPECT().Err().Return(nil)
					return splitter
				},
				unstructuredSerializerFn: func(t *testing.T) unstructuredSerializer {
					us := newMockUnstructuredSerializer(t)
					us.EXPECT().Decode(jobBytes, (*schema.GroupVersionKind)(nil), nil).Return(testJob, nil, nil)
					return us
				},
				unstructuredConverterFn: func(t *testing.T) unstructuredConverter {
					uc := newMockUnstructuredConverter(t)
					uc.EXPECT().ToUnstructured(testJob).Return(jobMap(), nil).Once()
					uc.EXPECT().ToUnstructured(testJobWithLabels).Return(jobMap(), nil).Once()
					uc.EXPECT().FromUnstructured(jobMap(), &batchv1.Job{}).Run(func(um map[string]interface{}, obj interface{}) {
						job := obj.(*batchv1.Job)
						*job = *testJob
					}).Return(nil)
					return uc
				},
				serializerFn: func(t *testing.T) genericYamlSerializer {
					ys := newMockGenericYamlSerializer(t)
					ys.EXPECT().Marshal(&unstructured.Unstructured{Object: jobMapWitLabels()}).Return(jobWithLabelsBytes, nil)
					return ys
				},
				labels: map[string]string{
					"k8s.cloudogu.com/component.name":    "k8s-blueprint-operator",
					"k8s.cloudogu.com/component.version": "1.2.3-4",
				},
			},
			renderedManifests:        bytes.NewBuffer(jobBytes),
			wantModifiedManifestsStr: fmt.Sprintf("%s\n---\n", jobWithLabelsBytes),
			wantErr:                  assert.NoError,
		},
		{
			name: "labels for pod-template in deployment",
			fields: fields{
				documentSplitterFn: func(t *testing.T) documentSplitter {
					return yamlutil.NewDocumentSplitter()
				},
				unstructuredSerializerFn: func(t *testing.T) unstructuredSerializer {
					return yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
				},
				unstructuredConverterFn: func(t *testing.T) unstructuredConverter {
					return runtime.DefaultUnstructuredConverter
				},
				serializerFn: func(t *testing.T) genericYamlSerializer {
					return yamlutil.NewSerializer()
				},
				labels: map[string]string{
					"k8s.cloudogu.com/component.name":    "k8s-test",
					"k8s.cloudogu.com/component.version": "1.2.3-4",
				},
			},
			renderedManifests:        bytes.NewBuffer(deploymentBytes),
			wantModifiedManifestsStr: deploymentWithLabelsString,
			wantErr:                  assert.NoError,
		},
		{
			name: "labels for pod-template in deployment with error",
			fields: fields{
				documentSplitterFn: func(t *testing.T) documentSplitter {
					return yamlutil.NewDocumentSplitter()
				},
				unstructuredSerializerFn: func(t *testing.T) unstructuredSerializer {
					return yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
				},
				unstructuredConverterFn: func(t *testing.T) unstructuredConverter {
					uc := newMockUnstructuredConverter(t)
					uc.EXPECT().ToUnstructured(mock.Anything).RunAndReturn(runtime.DefaultUnstructuredConverter.ToUnstructured)
					uc.EXPECT().FromUnstructured(mock.Anything, mock.Anything).Return(assert.AnError)
					return uc
				},
				serializerFn: func(t *testing.T) genericYamlSerializer {
					return yamlutil.NewSerializer()
				},
				labels: map[string]string{
					"k8s.cloudogu.com/component.name":    "k8s-test",
					"k8s.cloudogu.com/component.version": "1.2.3-4",
				},
			},
			renderedManifests:        bytes.NewBuffer(deploymentBytes),
			wantModifiedManifestsStr: "<nil>",
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, assert.AnError, i) &&
					assert.ErrorContains(t, err, "failed to add labels to Deployment: failed to convert resource to structured object:", i)
			},
		},
		{
			name: "labels for pod-template in statefulSet",
			fields: fields{
				documentSplitterFn: func(t *testing.T) documentSplitter {
					return yamlutil.NewDocumentSplitter()
				},
				unstructuredSerializerFn: func(t *testing.T) unstructuredSerializer {
					return yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
				},
				unstructuredConverterFn: func(t *testing.T) unstructuredConverter {
					return runtime.DefaultUnstructuredConverter
				},
				serializerFn: func(t *testing.T) genericYamlSerializer {
					return yamlutil.NewSerializer()
				},
				labels: map[string]string{
					"k8s.cloudogu.com/component.name":    "k8s-test",
					"k8s.cloudogu.com/component.version": "1.2.3-4",
				},
			},
			renderedManifests:        bytes.NewBuffer(statefulSetBytes),
			wantModifiedManifestsStr: statefulSetWithLabelsString,
			wantErr:                  assert.NoError,
		},
		{
			name: "labels for pod-template in statefulSet with error",
			fields: fields{
				documentSplitterFn: func(t *testing.T) documentSplitter {
					return yamlutil.NewDocumentSplitter()
				},
				unstructuredSerializerFn: func(t *testing.T) unstructuredSerializer {
					return yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
				},
				unstructuredConverterFn: func(t *testing.T) unstructuredConverter {
					uc := newMockUnstructuredConverter(t)
					uc.EXPECT().ToUnstructured(mock.Anything).RunAndReturn(runtime.DefaultUnstructuredConverter.ToUnstructured)
					uc.EXPECT().FromUnstructured(mock.Anything, mock.Anything).Return(assert.AnError)
					return uc
				},
				serializerFn: func(t *testing.T) genericYamlSerializer {
					return yamlutil.NewSerializer()
				},
				labels: map[string]string{
					"k8s.cloudogu.com/component.name":    "k8s-test",
					"k8s.cloudogu.com/component.version": "1.2.3-4",
				},
			},
			renderedManifests:        bytes.NewBuffer(statefulSetBytes),
			wantModifiedManifestsStr: "<nil>",
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, assert.AnError, i) &&
					assert.ErrorContains(t, err, "failed to add labels to StatefulSet: failed to convert resource to structured object:", i)
			},
		},
		{
			name: "labels for pod-template in daemonSet",
			fields: fields{
				documentSplitterFn: func(t *testing.T) documentSplitter {
					return yamlutil.NewDocumentSplitter()
				},
				unstructuredSerializerFn: func(t *testing.T) unstructuredSerializer {
					return yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
				},
				unstructuredConverterFn: func(t *testing.T) unstructuredConverter {
					return runtime.DefaultUnstructuredConverter
				},
				serializerFn: func(t *testing.T) genericYamlSerializer {
					return yamlutil.NewSerializer()
				},
				labels: map[string]string{
					"k8s.cloudogu.com/component.name":    "k8s-test",
					"k8s.cloudogu.com/component.version": "1.2.3-4",
				},
			},
			renderedManifests:        bytes.NewBuffer(daemonSetBytes),
			wantModifiedManifestsStr: daemonSetWithLabelsString,
			wantErr:                  assert.NoError,
		},
		{
			name: "labels for pod-template in daemonSet with error",
			fields: fields{
				documentSplitterFn: func(t *testing.T) documentSplitter {
					return yamlutil.NewDocumentSplitter()
				},
				unstructuredSerializerFn: func(t *testing.T) unstructuredSerializer {
					return yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
				},
				unstructuredConverterFn: func(t *testing.T) unstructuredConverter {
					uc := newMockUnstructuredConverter(t)
					uc.EXPECT().ToUnstructured(mock.Anything).RunAndReturn(runtime.DefaultUnstructuredConverter.ToUnstructured)
					uc.EXPECT().FromUnstructured(mock.Anything, mock.Anything).Return(assert.AnError)
					return uc
				},
				serializerFn: func(t *testing.T) genericYamlSerializer {
					return yamlutil.NewSerializer()
				},
				labels: map[string]string{
					"k8s.cloudogu.com/component.name":    "k8s-test",
					"k8s.cloudogu.com/component.version": "1.2.3-4",
				},
			},
			renderedManifests:        bytes.NewBuffer(daemonSetBytes),
			wantModifiedManifestsStr: "<nil>",
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, assert.AnError, i) &&
					assert.ErrorContains(t, err, "failed to add labels to DaemonSet: failed to convert resource to structured object:", i)
			},
		},
		{
			name: "labels for pod-template in job",
			fields: fields{
				documentSplitterFn: func(t *testing.T) documentSplitter {
					return yamlutil.NewDocumentSplitter()
				},
				unstructuredSerializerFn: func(t *testing.T) unstructuredSerializer {
					return yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
				},
				unstructuredConverterFn: func(t *testing.T) unstructuredConverter {
					return runtime.DefaultUnstructuredConverter
				},
				serializerFn: func(t *testing.T) genericYamlSerializer {
					return yamlutil.NewSerializer()
				},
				labels: map[string]string{
					"k8s.cloudogu.com/component.name":    "k8s-test",
					"k8s.cloudogu.com/component.version": "1.2.3-4",
				},
			},
			renderedManifests:        bytes.NewBuffer(jobBytes),
			wantModifiedManifestsStr: string(jobWithLabelsBytes),
			wantErr:                  assert.NoError,
		},
		{
			name: "labels for pod-template in job with error",
			fields: fields{
				documentSplitterFn: func(t *testing.T) documentSplitter {
					return yamlutil.NewDocumentSplitter()
				},
				unstructuredSerializerFn: func(t *testing.T) unstructuredSerializer {
					return yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
				},
				unstructuredConverterFn: func(t *testing.T) unstructuredConverter {
					uc := newMockUnstructuredConverter(t)
					uc.EXPECT().ToUnstructured(mock.Anything).RunAndReturn(runtime.DefaultUnstructuredConverter.ToUnstructured)
					uc.EXPECT().FromUnstructured(mock.Anything, mock.Anything).Return(assert.AnError)
					return uc
				},
				serializerFn: func(t *testing.T) genericYamlSerializer {
					return yamlutil.NewSerializer()
				},
				labels: map[string]string{
					"k8s.cloudogu.com/component.name":    "k8s-test",
					"k8s.cloudogu.com/component.version": "1.2.3-4",
				},
			},
			renderedManifests:        bytes.NewBuffer(jobBytes),
			wantModifiedManifestsStr: "<nil>",
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, assert.AnError, i) &&
					assert.ErrorContains(t, err, "failed to add labels to Job: failed to convert resource to structured object:", i)
			},
		},
		{
			name: "labels for pod-template in cronJob",
			fields: fields{
				documentSplitterFn: func(t *testing.T) documentSplitter {
					return yamlutil.NewDocumentSplitter()
				},
				unstructuredSerializerFn: func(t *testing.T) unstructuredSerializer {
					return yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
				},
				unstructuredConverterFn: func(t *testing.T) unstructuredConverter {
					return runtime.DefaultUnstructuredConverter
				},
				serializerFn: func(t *testing.T) genericYamlSerializer {
					return yamlutil.NewSerializer()
				},
				labels: map[string]string{
					"k8s.cloudogu.com/component.name":    "k8s-test",
					"k8s.cloudogu.com/component.version": "1.2.3-4",
				},
			},
			renderedManifests:        bytes.NewBuffer(cronJobBytes),
			wantModifiedManifestsStr: cronJobWithLabelsString,
			wantErr:                  assert.NoError,
		},
		{
			name: "labels for pod-template in cronJob with error",
			fields: fields{
				documentSplitterFn: func(t *testing.T) documentSplitter {
					return yamlutil.NewDocumentSplitter()
				},
				unstructuredSerializerFn: func(t *testing.T) unstructuredSerializer {
					return yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
				},
				unstructuredConverterFn: func(t *testing.T) unstructuredConverter {
					uc := newMockUnstructuredConverter(t)
					uc.EXPECT().ToUnstructured(mock.Anything).RunAndReturn(runtime.DefaultUnstructuredConverter.ToUnstructured)
					uc.EXPECT().FromUnstructured(mock.Anything, mock.Anything).Return(assert.AnError)
					return uc
				},
				serializerFn: func(t *testing.T) genericYamlSerializer {
					return yamlutil.NewSerializer()
				},
				labels: map[string]string{
					"k8s.cloudogu.com/component.name":    "k8s-test",
					"k8s.cloudogu.com/component.version": "1.2.3-4",
				},
			},
			renderedManifests:        bytes.NewBuffer(cronJobBytes),
			wantModifiedManifestsStr: "<nil>",
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, assert.AnError, i) &&
					assert.ErrorContains(t, err, "failed to add labels to CronJob: failed to convert resource to structured object:", i)
			},
		},
		{
			name: "test integration dogu-operator",
			fields: fields{
				documentSplitterFn: func(t *testing.T) documentSplitter {
					return yamlutil.NewDocumentSplitter()
				},
				unstructuredSerializerFn: func(t *testing.T) unstructuredSerializer {
					return yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
				},
				unstructuredConverterFn: func(t *testing.T) unstructuredConverter {
					return runtime.DefaultUnstructuredConverter
				},
				serializerFn: func(t *testing.T) genericYamlSerializer {
					return yamlutil.NewSerializer()
				},
				labels: map[string]string{
					"k8s.cloudogu.com/component.name":    "k8s-dogu-operator",
					"k8s.cloudogu.com/component.version": "1.2.3-4",
				},
			},
			renderedManifests:        bytes.NewBuffer(doguOpBytes),
			wantModifiedManifestsStr: doguOpWithLabelsStr,
			wantErr:                  assert.NoError,
		},
		{
			name: "test integration longhorn",
			fields: fields{
				documentSplitterFn: func(t *testing.T) documentSplitter {
					return yamlutil.NewDocumentSplitter()
				},
				unstructuredSerializerFn: func(t *testing.T) unstructuredSerializer {
					return yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
				},
				unstructuredConverterFn: func(t *testing.T) unstructuredConverter {
					return runtime.DefaultUnstructuredConverter
				},
				serializerFn: func(t *testing.T) genericYamlSerializer {
					return yamlutil.NewSerializer()
				},
				labels: map[string]string{
					"k8s.cloudogu.com/component.name":    "k8s-longhorn",
					"k8s.cloudogu.com/component.version": "1.5.1-4",
				},
			},
			renderedManifests:        bytes.NewBuffer(longhornBytes),
			wantModifiedManifestsStr: longhornWithLabelsStr,
			wantErr:                  assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &PostRenderer{
				documentSplitter:       tt.fields.documentSplitterFn(t),
				unstructuredSerializer: tt.fields.unstructuredSerializerFn(t),
				unstructuredConverter:  tt.fields.unstructuredConverterFn(t),
				serializer:             tt.fields.serializerFn(t),
				labels:                 tt.fields.labels,
			}
			gotModifiedManifests, err := c.Run(tt.renderedManifests)
			tt.wantErr(t, err)
			assert.Equal(t, tt.wantModifiedManifestsStr, gotModifiedManifests.String())
		})
	}
}

func jobMap() map[string]interface{} {
	return map[string]interface{}{"apiVersion": "batch/v1", "kind": "Job", "metadata": map[string]interface{}{"name": "hello"},
		"spec": map[string]interface{}{"template": map[string]interface{}{"spec": map[string]interface{}{
			"containers": []map[string]interface{}{{"name": "hello", "image": "busybox", "command": []interface{}{"sh", "-c", `echo "Hello, Kubernetes!" && sleep 3600`}}},
		}}},
	}
}

func jobMapWitLabels() map[string]interface{} {
	return map[string]interface{}{"apiVersion": "batch/v1", "kind": "Job",
		"metadata": map[string]interface{}{"name": "hello", "labels": map[string]interface{}{
			"k8s.cloudogu.com/component.name":    "k8s-blueprint-operator",
			"k8s.cloudogu.com/component.version": "1.2.3-4",
		}},
		"spec": map[string]interface{}{"template": map[string]interface{}{"spec": map[string]interface{}{
			"containers": []map[string]interface{}{{"name": "hello", "image": "busybox", "command": []interface{}{"sh", "-c", `echo "Hello, Kubernetes!" && sleep 3600`}}},
		}}},
	}
}

func TestNewPostRenderer(t *testing.T) {
	renderer := NewPostRenderer(map[string]string{
		"k8s.cloudogu.com/component.name":    "k8s-blueprint-operator",
		"k8s.cloudogu.com/component.version": "1.2.3-4",
	})
	assert.NotEmpty(t, renderer)
}

func Test_addLabelsToStructured(t *testing.T) {
	t.Run("should fail to convert to unstructured", func(t *testing.T) {
		mockConverter := newMockUnstructuredConverter(t)
		mockConverter.EXPECT().FromUnstructured(mock.Anything, &appsv1.DaemonSet{}).Return(nil)
		mockConverter.EXPECT().ToUnstructured(mock.Anything).Return(nil, assert.AnError)

		c := &PostRenderer{
			unstructuredConverter: mockConverter,
		}
		_, err := addLabelsToStructured(c, nil, &appsv1.DaemonSet{}, func(a *appsv1.DaemonSet) objectWithLabels {
			return &a.Spec.Template
		})

		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to convert resource to unstructured object:")
	})
}
