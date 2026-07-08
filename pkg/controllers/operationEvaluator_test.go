package controllers

import (
	"testing"
	"time"

	k8sv1 "github.com/cloudogu/k8s-component-lib/api/v1"
	"github.com/cloudogu/k8s-component-operator/pkg/yaml"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/release"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Test_defaultOperationEvaluatorFactory_NewOperationEvaluator(t *testing.T) {
	// given
	helmClientMock := newMockHelmClient(t)
	recorderMock := newMockEventRecorder(t)
	yamlSerializer := yaml.NewSerializer()
	readerMock := newMockConfigMapRefReader(t)
	timeout := 5 * time.Minute

	sut := defaultOperationEvaluatorFactory{
		recorder:       recorderMock,
		timeout:        timeout,
		yamlSerializer: yamlSerializer,
		reader:         readerMock,
	}

	// when
	actual := sut.NewOperationEvaluator(helmClientMock)

	// then
	evaluator, ok := actual.(*defaultOperationEvaluator)
	require.True(t, ok)
	assert.Same(t, helmClientMock, evaluator.helmClient)
	assert.Same(t, recorderMock, evaluator.recorder)
	assert.Equal(t, timeout, evaluator.timeout)
	assert.Same(t, yamlSerializer, evaluator.yamlSerializer)
	assert.Same(t, readerMock, evaluator.reader)
}

func Test_defaultOperationEvaluator_getChangeOperation(t *testing.T) {
	t.Run("should fail on error getting helm releases", func(t *testing.T) {
		// given
		component := getComponent("ecosystem", "k8s", "", "dogu-op", "0.1.0")
		mockHelmClient := newMockHelmClient(t)
		mockHelmClient.EXPECT().ListDeployedReleases().Return(nil, assert.AnError)

		sut := defaultOperationEvaluator{helmClient: mockHelmClient}

		// when
		_, err := sut.getChangeOperation(testCtx, component)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to get deployed helm releases")
	})

	t.Run("should fail on error parsing component version", func(t *testing.T) {
		// given
		component := getComponent("ecosystem", "k8s", "", "dogu-op", "notvalidsemver")
		mockHelmClient := newMockHelmClient(t)
		helmReleases := []*release.Release{{Name: "dogu-op", Namespace: "ecosystem", Chart: &chart.Chart{Metadata: &chart.Metadata{AppVersion: "0.0.1"}}}}
		mockHelmClient.EXPECT().ListDeployedReleases().Return(helmReleases, nil)

		sut := defaultOperationEvaluator{helmClient: mockHelmClient}

		// when
		_, err := sut.getChangeOperation(testCtx, component)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "failed to parse component version")
	})

	t.Run("should fail on error getting release values", func(t *testing.T) {
		// given
		component := getComponent("ecosystem", "k8s", "", "dogu-op", "0.0.1")
		mockHelmClient := newMockHelmClient(t)
		helmReleases := []*release.Release{{Name: "dogu-op", Namespace: "ecosystem", Chart: &chart.Chart{Metadata: &chart.Metadata{AppVersion: "0.0.1"}}}}
		mockHelmClient.EXPECT().ListDeployedReleases().Return(helmReleases, nil)
		mockHelmClient.EXPECT().GetReleaseValues("dogu-op", false).Return(nil, assert.AnError)

		sut := defaultOperationEvaluator{helmClient: mockHelmClient}

		// when
		_, err := sut.getChangeOperation(testCtx, component)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to compare Values.yaml files of component")
		assert.ErrorContains(t, err, "failed to get values.yaml from release")
	})

	t.Run("should fail on error getting component values", func(t *testing.T) {
		// given
		component := getComponent("ecosystem", "k8s", "", "dogu-op", "0.0.1")
		component.Spec.ValuesConfigRef = &k8sv1.Reference{}
		mockHelmClient := newMockHelmClient(t)
		helmReleases := []*release.Release{{Name: "dogu-op", Namespace: "ecosystem", Chart: &chart.Chart{Metadata: &chart.Metadata{AppVersion: "0.0.1"}}}}
		mockHelmClient.EXPECT().ListDeployedReleases().Return(helmReleases, nil)
		mockHelmClient.EXPECT().GetReleaseValues("dogu-op", false).Return(map[string]interface{}{}, nil)
		mockHelmClient.EXPECT().GetChartSpecValues(mock.Anything).Return(nil, assert.AnError)
		configMapRefReaderMock := newMockConfigMapRefReader(t)
		configMapRefReaderMock.EXPECT().GetValues(testCtx, &k8sv1.Reference{}).Return("", nil)

		sut := defaultOperationEvaluator{
			helmClient:     mockHelmClient,
			timeout:        defaultHelmClientTimeoutMins,
			yamlSerializer: yaml.NewSerializer(),
			reader:         configMapRefReaderMock,
		}

		// when
		_, err := sut.getChangeOperation(testCtx, component)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to compare Values.yaml files of component")
		assert.ErrorContains(t, err, "failed to get values.yaml from component")
	})

	t.Run("should return downgrade-operation on downgrade", func(t *testing.T) {
		// given
		component := getComponent("ecosystem", "k8s", "", "dogu-op", "0.0.0")
		mockHelmClient := newMockHelmClient(t)
		helmReleases := []*release.Release{{Name: "dogu-op", Namespace: "ecosystem", Chart: &chart.Chart{Metadata: &chart.Metadata{AppVersion: "0.0.1"}}}}
		mockHelmClient.EXPECT().ListDeployedReleases().Return(helmReleases, nil)

		sut := defaultOperationEvaluator{helmClient: mockHelmClient}

		// when
		op, err := sut.getChangeOperation(testCtx, component)

		// then
		require.NoError(t, err)
		assert.Equal(t, Downgrade, op)
	})

	t.Run("should return upgrade-operation on upgrade if deploy namespace is set", func(t *testing.T) {
		// given
		component := getComponent("ecosystem", "k8s", "deploy-namespace", "dogu-op", "0.0.1-2")
		mockHelmClient := newMockHelmClient(t)
		helmReleases := []*release.Release{{Name: "dogu-op", Namespace: "deploy-namespace", Chart: &chart.Chart{Metadata: &chart.Metadata{AppVersion: "0.0.1-1"}}}}
		mockHelmClient.EXPECT().ListDeployedReleases().Return(helmReleases, nil)

		sut := defaultOperationEvaluator{helmClient: mockHelmClient}

		// when
		op, err := sut.getChangeOperation(testCtx, component)

		// then
		require.NoError(t, err)
		assert.Equal(t, Upgrade, op)
	})

	t.Run("should return error if deploy namespace is not the same as release namespace", func(t *testing.T) {
		// given
		component := getComponent("ecosystem", "k8s", "deploy-namespace", "dogu-op", "0.0.1-2")
		mockHelmClient := newMockHelmClient(t)
		mockRecorder := newMockEventRecorder(t)
		mockRecorder.EXPECT().Eventf(component, corev1.EventTypeWarning, UpgradeEventReason, "Deploy namespace mismatch (CR: %q; deployed: %q). Deploy namespace declaration is only allowed on install. Revert deploy namespace change to prevent failing upgrade.", "deploy-namespace", "ecosystem").Return()
		helmReleases := []*release.Release{{Name: "dogu-op", Namespace: "ecosystem", Chart: &chart.Chart{Metadata: &chart.Metadata{AppVersion: "0.0.1-2"}}}}
		mockHelmClient.EXPECT().ListDeployedReleases().Return(helmReleases, nil)

		sut := defaultOperationEvaluator{
			helmClient: mockHelmClient,
			recorder:   mockRecorder,
		}

		// when
		_, err := sut.getChangeOperation(testCtx, component)

		// then
		require.Error(t, err)
	})

	t.Run("should return upgrade-operation on upgrade", func(t *testing.T) {
		// given
		component := getComponent("ecosystem", "k8s", "", "dogu-op", "0.0.2")
		mockHelmClient := newMockHelmClient(t)
		helmReleases := []*release.Release{{Name: "dogu-op", Namespace: "ecosystem", Chart: &chart.Chart{Metadata: &chart.Metadata{AppVersion: "0.0.1"}}}}
		mockHelmClient.EXPECT().ListDeployedReleases().Return(helmReleases, nil)

		sut := defaultOperationEvaluator{helmClient: mockHelmClient}

		// when
		op, err := sut.getChangeOperation(testCtx, component)

		// then
		require.NoError(t, err)
		assert.Equal(t, Upgrade, op)
	})

	t.Run("should return upgrade-operation on same version, but values-yaml difference", func(t *testing.T) {
		// given
		component := getComponent("ecosystem", "k8s", "", "dogu-op", "0.0.2")
		component.Spec.ValuesConfigRef = &k8sv1.Reference{}
		mockHelmClient := newMockHelmClient(t)
		helmReleases := []*release.Release{{Name: "dogu-op", Namespace: "ecosystem", Chart: &chart.Chart{Metadata: &chart.Metadata{AppVersion: "0.0.2"}}}}
		mockHelmClient.EXPECT().ListDeployedReleases().Return(helmReleases, nil)
		mockHelmClient.EXPECT().GetReleaseValues("dogu-op", false).Return(map[string]interface{}{"foo": "bar", "baz": "buz"}, nil)
		mockHelmClient.EXPECT().GetChartSpecValues(mock.Anything).Return(map[string]interface{}{"foo": "bar", "baz": "xyz"}, nil)
		configMapRefReaderMock := newMockConfigMapRefReader(t)
		configMapRefReaderMock.EXPECT().GetValues(testCtx, &k8sv1.Reference{}).Return("", nil)

		sut := defaultOperationEvaluator{
			helmClient:     mockHelmClient,
			timeout:        defaultHelmClientTimeoutMins,
			yamlSerializer: yaml.NewSerializer(),
			reader:         configMapRefReaderMock,
		}

		// when
		op, err := sut.getChangeOperation(testCtx, component)

		// then
		require.NoError(t, err)
		assert.Equal(t, Upgrade, op)
	})

	t.Run("should return ignore-operation on same version and same values-yaml values", func(t *testing.T) {
		// given
		component := getComponent("ecosystem", "k8s", "", "dogu-op", "0.0.2")
		component.Spec.ValuesConfigRef = &k8sv1.Reference{}
		mockHelmClient := newMockHelmClient(t)
		helmReleases := []*release.Release{{Name: "dogu-op", Namespace: "ecosystem", Chart: &chart.Chart{Metadata: &chart.Metadata{AppVersion: "0.0.2"}}}}
		mockHelmClient.EXPECT().ListDeployedReleases().Return(helmReleases, nil)
		mockHelmClient.EXPECT().GetReleaseValues("dogu-op", false).Return(map[string]interface{}{"foo": "bar", "baz": "buz"}, nil)
		mockHelmClient.EXPECT().GetChartSpecValues(mock.Anything).Return(map[string]interface{}{"foo": "bar", "baz": "buz"}, nil)
		configMapRefReaderMock := newMockConfigMapRefReader(t)
		configMapRefReaderMock.EXPECT().GetValues(testCtx, &k8sv1.Reference{}).Return("", nil)

		sut := defaultOperationEvaluator{
			helmClient:     mockHelmClient,
			timeout:        defaultHelmClientTimeoutMins,
			yamlSerializer: yaml.NewSerializer(),
			reader:         configMapRefReaderMock,
		}

		// when
		op, err := sut.getChangeOperation(testCtx, component)

		// then
		require.NoError(t, err)
		assert.Equal(t, Ignore, op)
	})

	t.Run("should return ignore-operation on same version and different zero-length maps", func(t *testing.T) {
		// given
		component := getComponent("ecosystem", "k8s", "", "dogu-op", "0.0.2")
		component.Spec.ValuesConfigRef = &k8sv1.Reference{}
		mockHelmClient := newMockHelmClient(t)
		helmReleases := []*release.Release{{Name: "dogu-op", Namespace: "ecosystem", Chart: &chart.Chart{Metadata: &chart.Metadata{AppVersion: "0.0.2"}}}}
		mockHelmClient.EXPECT().ListDeployedReleases().Return(helmReleases, nil)
		mockHelmClient.EXPECT().GetReleaseValues("dogu-op", false).Return(map[string]interface{}(nil), nil)
		mockHelmClient.EXPECT().GetChartSpecValues(mock.Anything).Return(map[string]interface{}{}, nil)
		configMapRefReaderMock := newMockConfigMapRefReader(t)
		configMapRefReaderMock.EXPECT().GetValues(testCtx, &k8sv1.Reference{}).Return("", nil)

		sut := defaultOperationEvaluator{
			helmClient:     mockHelmClient,
			timeout:        defaultHelmClientTimeoutMins,
			yamlSerializer: yaml.NewSerializer(),
			reader:         configMapRefReaderMock,
		}

		// when
		op, err := sut.getChangeOperation(testCtx, component)

		// then
		require.NoError(t, err)
		assert.Equal(t, Ignore, op)
	})

	t.Run("should return ignore-operation on same version", func(t *testing.T) {
		// given
		component := getComponent("ecosystem", "k8s", "", "dogu-op", "0.0.1")
		component.Spec.ValuesConfigRef = &k8sv1.Reference{}
		mockHelmClient := newMockHelmClient(t)
		helmReleases := []*release.Release{{Name: "dogu-op", Namespace: "ecosystem", Chart: &chart.Chart{Metadata: &chart.Metadata{AppVersion: "0.0.1"}}}}
		mockHelmClient.EXPECT().ListDeployedReleases().Return(helmReleases, nil)
		mockHelmClient.EXPECT().GetReleaseValues("dogu-op", false).Return(map[string]interface{}{}, nil)
		mockHelmClient.EXPECT().GetChartSpecValues(mock.Anything).Return(map[string]interface{}{}, nil)
		configMapRefReaderMock := newMockConfigMapRefReader(t)
		configMapRefReaderMock.EXPECT().GetValues(testCtx, &k8sv1.Reference{}).Return("", nil)

		sut := defaultOperationEvaluator{
			helmClient:     mockHelmClient,
			timeout:        defaultHelmClientTimeoutMins,
			yamlSerializer: yaml.NewSerializer(),
			reader:         configMapRefReaderMock,
		}

		// when
		op, err := sut.getChangeOperation(testCtx, component)

		// then
		require.NoError(t, err)
		assert.Equal(t, Ignore, op)
	})

	t.Run("should return ignore-operation when no release is found", func(t *testing.T) {
		// given
		component := getComponent("ecosystem", "k8s", "", "dogu-op", "0.0.1")
		mockHelmClient := newMockHelmClient(t)
		var helmReleases []*release.Release
		mockHelmClient.EXPECT().ListDeployedReleases().Return(helmReleases, nil)

		sut := defaultOperationEvaluator{helmClient: mockHelmClient}

		// when
		op, err := sut.getChangeOperation(testCtx, component)

		// then
		require.NoError(t, err)
		assert.Equal(t, Ignore, op)
	})
}

func Test_defaultOperationEvaluator_EvaluateRequiredOperation(t *testing.T) {
	t.Run("should return install on status installing", func(t *testing.T) {
		// given
		component := getComponent("ecosystem", "k8s", "", "dogu-op", "0.0.0")
		component.Status.Status = "installing"
		sut := defaultOperationEvaluator{}

		// when
		requiredOperation, err := sut.EvaluateRequiredOperation(testCtx, component)

		// then
		require.NoError(t, err)
		assert.Equal(t, Install, requiredOperation)
	})

	t.Run("should return delete on status deleting", func(t *testing.T) {
		// given
		component := getComponent("ecosystem", "k8s", "", "dogu-op", "0.0.0")
		component.Status.Status = "deleting"
		sut := defaultOperationEvaluator{}

		// when
		requiredOperation, err := sut.EvaluateRequiredOperation(testCtx, component)

		// then
		require.NoError(t, err)
		assert.Equal(t, Delete, requiredOperation)
	})

	t.Run("should return upgrade on status upgrading", func(t *testing.T) {
		// given
		component := getComponent("ecosystem", "k8s", "", "dogu-op", "0.0.0")
		component.Status.Status = "upgrading"
		sut := defaultOperationEvaluator{}

		// when
		requiredOperation, err := sut.EvaluateRequiredOperation(testCtx, component)

		// then
		require.NoError(t, err)
		assert.Equal(t, Upgrade, requiredOperation)
	})

	t.Run("should return ignore on unrecognized status", func(t *testing.T) {
		// given
		component := getComponent("ecosystem", "k8s", "", "dogu-op", "0.0.0")
		component.Status.Status = "foobar"
		sut := defaultOperationEvaluator{}

		// when
		requiredOperation, err := sut.EvaluateRequiredOperation(testCtx, component)

		// then
		require.NoError(t, err)
		assert.Equal(t, Ignore, requiredOperation)
	})

	t.Run("should return install on tryToInstall status", func(t *testing.T) {
		// given
		component := getComponent("ecosystem", "k8s", "", "dogu-op", "0.0.0")
		component.Status.Status = "tryToInstall"
		sut := defaultOperationEvaluator{}

		// when
		requiredOperation, err := sut.EvaluateRequiredOperation(testCtx, component)

		// then
		require.NoError(t, err)
		assert.Equal(t, Install, requiredOperation)
	})

	t.Run("should return upgrade on tryToUpgrade status", func(t *testing.T) {
		// given
		componentName := "dogu-op"
		component := getComponent("ecosystem", "k8s", "", componentName, "0.0.1")
		component.Status.Status = "tryToUpgrade"
		helmMock := newMockHelmClient(t)
		installedReleases := []*release.Release{{Namespace: "ecosystem", Name: componentName, Chart: &chart.Chart{Metadata: &chart.Metadata{AppVersion: "0.0.0"}}}}
		helmMock.EXPECT().ListDeployedReleases().Return(installedReleases, nil)
		sut := defaultOperationEvaluator{helmClient: helmMock}

		// when
		requiredOperation, err := sut.EvaluateRequiredOperation(testCtx, component)

		// then
		require.NoError(t, err)
		assert.Equal(t, Upgrade, requiredOperation)
	})

	t.Run("should return delete on tryToDelete status", func(t *testing.T) {
		// given
		componentName := "dogu-op"
		component := getComponent("ecosystem", "k8s", "", componentName, "0.0.0")
		component.Status.Status = "tryToDelete"
		timeNow := v1.NewTime(time.Now())
		component.DeletionTimestamp = &timeNow
		sut := defaultOperationEvaluator{}

		// when
		requiredOperation, err := sut.EvaluateRequiredOperation(testCtx, component)

		// then
		require.NoError(t, err)
		assert.Equal(t, Delete, requiredOperation)
	})
}
