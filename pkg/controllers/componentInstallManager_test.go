package controllers

//
// import (
// 	"context"
// 	"errors"
// 	"github.com/cloudogu/k8s-dogu-operator/internal"
// 	"github.com/cloudogu/k8s-dogu-operator/internal/mocks/external"
// 	"testing"
//
// 	apierrors "k8s.io/apimachinery/pkg/api/errors"
//
// 	imagev1 "github.com/google/go-containerregistry/pkg/v1"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/mock"
// 	"github.com/stretchr/testify/require"
// 	corev1 "k8s.io/api/core/v1"
// 	"k8s.io/apimachinery/pkg/runtime"
// 	"k8s.io/clientset-go/rest"
// 	"k8s.io/clientset-go/tools/clientcmd/api"
// 	ctrl "sigs.k8s.io/controller-runtime"
// 	"sigs.k8s.io/controller-runtime/pkg/clientset"
// 	"sigs.k8s.io/controller-runtime/pkg/clientset/fake"
//
// 	"github.com/cloudogu/cesapp-lib/core"
// 	cesmocks "github.com/cloudogu/cesapp-lib/registry/mocks"
//
// 	k8sv1 "github.com/cloudogu/k8s-dogu-operator/api/v1"
// 	"github.com/cloudogu/k8s-dogu-operator/controllers/config"
// 	"github.com/cloudogu/k8s-dogu-operator/controllers/resource"
// 	"github.com/cloudogu/k8s-dogu-operator/internal/mocks"
// )
//
// type doguInstallManagerWithMocks struct {
// 	installManager            *componentInstallManager
// 	localDoguFetcher          *mocks.LocalDoguFetcher
// 	resourceDoguFetcher       *mocks.ResourceDoguFetcher
// 	imageRegistryMock         *mocks.ImageRegistry
// 	doguRegistratorMock       *mocks.DoguRegistrator
// 	dependencyValidatorMock   *mocks.DependencyValidator
// 	serviceAccountCreatorMock *mocks.ServiceAccountCreator
// 	doguSecretHandlerMock     *mocks.DoguSecretHandler
// 	applierMock               *mocks.Applier
// 	fileExtractorMock         *mocks.FileExtractor
// 	clientset                    clientset.WithWatch
// 	resourceUpserter          *mocks.ResourceUpserter
// 	recorder                  *external.EventRecorder
// 	execPodFactory            *mocks.ExecPodFactory
// }
//
// func (d *doguInstallManagerWithMocks) AssertMocks(t *testing.T) {
// 	t.Helper()
// 	mock.AssertExpectationsForObjects(t,
// 		d.imageRegistryMock,
// 		d.doguRegistratorMock,
// 		d.dependencyValidatorMock,
// 		d.serviceAccountCreatorMock,
// 		d.doguSecretHandlerMock,
// 		d.applierMock,
// 		d.fileExtractorMock,
// 		d.localDoguFetcher,
// 		d.resourceDoguFetcher,
// 		d.recorder,
// 		d.resourceUpserter,
// 		d.execPodFactory,
// 	)
// }
//
// func getDoguInstallManagerWithMocks(t *testing.T, scheme *runtime.Scheme) doguInstallManagerWithMocks {
// 	k8sClient := fake.NewClientBuilder().WithScheme(scheme).Build()
// 	limitPatcher := &mocks.LimitPatcher{}
// 	limitPatcher.On("RetrievePodLimits", mock.Anything).Return(mocks.NewDoguLimits(t), nil)
// 	limitPatcher.On("PatchDeployment", mock.Anything, mock.Anything).Return(nil)
// 	upserter := &mocks.ResourceUpserter{}
// 	imageRegistry := &mocks.ImageRegistry{}
// 	doguRegistrator := &mocks.DoguRegistrator{}
// 	dependencyValidator := &mocks.DependencyValidator{}
// 	serviceAccountCreator := &mocks.ServiceAccountCreator{}
// 	doguSecretHandler := &mocks.DoguSecretHandler{}
// 	mockedApplier := &mocks.Applier{}
// 	fileExtract := mocks.NewFileExtractor(t)
// 	eventRecorderMock := external.NewEventRecorder(t)
// 	localDoguFetcher := mocks.NewLocalDoguFetcher(t)
// 	resourceDoguFetcher := mocks.NewResourceDoguFetcher(t)
// 	collectApplier := resource.NewCollectApplier(mockedApplier)
// 	podFactory := mocks.NewExecPodFactory(t)
//
// 	componentInstallManager := &componentInstallManager{
// 		clientset:                k8sClient,
// 		recorder:              eventRecorderMock,
// 		imageRegistry:         imageRegistry,
// 		doguRegistrator:       doguRegistrator,
// 		localDoguFetcher:      localDoguFetcher,
// 		resourceDoguFetcher:   resourceDoguFetcher,
// 		dependencyValidator:   dependencyValidator,
// 		serviceAccountCreator: serviceAccountCreator,
// 		doguSecretHandler:     doguSecretHandler,
// 		fileExtractor:         fileExtract,
// 		collectApplier:        collectApplier,
// 		resourceUpserter:      upserter,
// 		execPodFactory:        podFactory,
// 	}
//
// 	return doguInstallManagerWithMocks{
// 		installManager:            componentInstallManager,
// 		clientset:                    k8sClient,
// 		recorder:                  eventRecorderMock,
// 		localDoguFetcher:          localDoguFetcher,
// 		resourceDoguFetcher:       resourceDoguFetcher,
// 		imageRegistryMock:         imageRegistry,
// 		doguRegistratorMock:       doguRegistrator,
// 		dependencyValidatorMock:   dependencyValidator,
// 		serviceAccountCreatorMock: serviceAccountCreator,
// 		doguSecretHandlerMock:     doguSecretHandler,
// 		fileExtractorMock:         fileExtract,
// 		applierMock:               mockedApplier,
// 		resourceUpserter:          upserter,
// 		execPodFactory:            podFactory,
// 	}
// }
//
// func getDoguInstallManagerTestData(t *testing.T) (*k8sv1.Dogu, *core.Dogu, *corev1.ConfigMap, *imagev1.ConfigFile) {
// 	ldapCr := readDoguCr(t, ldapCrBytes)
// 	ldapDogu := readDoguDescriptor(t, ldapDoguDescriptorBytes)
// 	ldapDoguDescriptor := readDoguDevelopmentMap(t, ldapDoguDevelopmentMapBytes)
// 	imageConfig := readImageConfig(t, imageConfigBytes)
// 	return ldapCr, ldapDogu, ldapDoguDescriptor.ToConfigMap(), imageConfig
// }
//
// func TestNewDoguInstallManager(t *testing.T) {
// 	// override default controller method to retrieve a kube config
// 	oldGetConfigOrDieDelegate := ctrl.GetConfigOrDie
// 	defer func() { ctrl.GetConfigOrDie = oldGetConfigOrDieDelegate }()
// 	ctrl.GetConfigOrDie = func() *rest.Config {
// 		return &rest.Config{}
// 	}
//
// 	t.Run("success", func(t *testing.T) {
// 		// given
// 		myClient := fake.NewClientBuilder().WithScheme(runtime.NewScheme()).Build()
// 		operatorConfig := &config.OperatorConfig{}
// 		operatorConfig.Namespace = "test"
// 		cesRegistry := &cesmocks.Registry{}
// 		doguRegistry := &cesmocks.DoguRegistry{}
// 		eventRecorder := &external.EventRecorder{}
// 		cesRegistry.On("DoguRegistry").Return(doguRegistry)
//
// 		// when
// 		doguManager, err := NewComponentInstallManager(myClient, operatorConfig, cesRegistry, eventRecorder)
//
// 		// then
// 		require.NoError(t, err)
// 		require.NotNil(t, doguManager)
// 		mock.AssertExpectationsForObjects(t, cesRegistry, doguRegistry)
// 	})
//
// 	t.Run("fail when creating clientset", func(t *testing.T) {
// 		// given
//
// 		// override default controller method to return a config that fail the clientset creation
// 		oldGetConfigOrDieDelegate := ctrl.GetConfigOrDie
// 		defer func() { ctrl.GetConfigOrDie = oldGetConfigOrDieDelegate }()
// 		ctrl.GetConfigOrDie = func() *rest.Config {
// 			return &rest.Config{ExecProvider: &api.ExecConfig{}, AuthProvider: &api.AuthProviderConfig{}}
// 		}
//
// 		myClient := fake.NewClientBuilder().WithScheme(runtime.NewScheme()).Build()
// 		operatorConfig := &config.OperatorConfig{}
// 		operatorConfig.Namespace = "test"
// 		cesRegistry := &cesmocks.Registry{}
// 		eventRecorder := &external.EventRecorder{}
//
// 		// when
// 		doguManager, err := NewComponentInstallManager(myClient, operatorConfig, cesRegistry, eventRecorder)
//
// 		// then
// 		require.Error(t, err)
// 		require.Nil(t, doguManager)
// 	})
// }
//
// func Test_doguInstallManager_Install(t *testing.T) {
// 	ctx := context.Background()
//
// 	t.Run("successfully install a dogu", func(t *testing.T) {
// 		// given
// 		managerWithMocks := getDoguInstallManagerWithMocks(t, getTestScheme())
// 		ldapCr, ldapDogu, _, imageConfig := getDoguInstallManagerTestData(t)
//
// 		managerWithMocks.resourceDoguFetcher.On("FetchWithResource", ctx, ldapCr).Return(ldapDogu, nil, nil)
// 		managerWithMocks.imageRegistryMock.On("PullImageConfig", mock.Anything, mock.Anything).Return(imageConfig, nil)
// 		managerWithMocks.doguRegistratorMock.On("RegisterNewDogu", mock.Anything, mock.Anything, mock.Anything).Return(nil)
// 		managerWithMocks.dependencyValidatorMock.On("ValidateDependencies", ctx, mock.Anything).Return(nil)
// 		managerWithMocks.serviceAccountCreatorMock.On("CreateAll", mock.Anything, mock.Anything, mock.Anything).Return(nil)
// 		managerWithMocks.doguSecretHandlerMock.On("WriteDoguSecretsToRegistry", mock.Anything, mock.Anything).Return(nil)
//
// 		yamlResult := map[string]string{"my-custom-resource.yml": "kind: Namespace"}
// 		managerWithMocks.fileExtractorMock.On("ExtractK8sResourcesFromContainer", mock.Anything, mock.Anything, mock.Anything).Return(yamlResult, nil)
// 		_ = managerWithMocks.installManager.clientset.Create(ctx, ldapCr)
//
// 		managerWithMocks.applierMock.On("ApplyWithOwner", mock.Anything, "", ldapCr).Return(nil)
// 		managerWithMocks.resourceUpserter.On("UpsertDoguDeployment", ctx, ldapCr, ldapDogu, mock.Anything).Once().Return(nil, nil)
// 		managerWithMocks.resourceUpserter.On("UpsertDoguService", ctx, ldapCr, imageConfig).Once().Return(nil, nil)
// 		managerWithMocks.resourceUpserter.On("UpsertDoguExposedServices", ctx, ldapCr, ldapDogu).Once().Return(nil, nil)
// 		managerWithMocks.resourceUpserter.On("UpsertDoguPVCs", ctx, ldapCr, ldapDogu).Once().Return(nil, nil)
//
// 		managerWithMocks.recorder.On("Event", mock.Anything, corev1.EventTypeNormal, InstallEventReason, "Checking dependencies...").
// 			On("Event", mock.Anything, corev1.EventTypeNormal, InstallEventReason, "Registering in the local dogu registry...").
// 			On("Event", mock.Anything, corev1.EventTypeNormal, InstallEventReason, "Creating required service accounts...").
// 			On("Eventf", mock.Anything, corev1.EventTypeNormal, InstallEventReason, "Pulling dogu image %s...", "registry.cloudogu.com/official/ldap:2.4.48-4").
// 			On("Eventf", mock.Anything, corev1.EventTypeNormal, InstallEventReason, "Starting execPod...").
// 			On("Eventf", mock.Anything, corev1.EventTypeNormal, InstallEventReason, "Creating custom dogu resources to the cluster: [%s]", "my-custom-resource.yml").
// 			On("Event", mock.Anything, corev1.EventTypeNormal, InstallEventReason, "Creating kubernetes resources...")
// 		execPod := mocks.NewExecPod(t)
// 		execPod.On("Create", testCtx).Return(nil)
// 		execPod.On("Delete", testCtx).Return(nil)
// 		managerWithMocks.execPodFactory.On("NewExecPod", internal.VolumeModeInstall, ldapCr, ldapDogu, mock.Anything).Return(execPod, nil)
//
// 		// when
// 		err := managerWithMocks.installManager.Install(ctx, ldapCr)
//
// 		// then
// 		require.NoError(t, err)
// 		managerWithMocks.AssertMocks(t)
// 	})
//
// 	t.Run("successfully install dogu with custom descriptor", func(t *testing.T) {
// 		// given
// 		managerWithMocks := getDoguInstallManagerWithMocks(t, getTestScheme())
// 		ldapCr, ldapDogu, ldapDevelopmentDoguMap, imageConfig := getDoguInstallManagerTestData(t)
// 		developmentDoguMap := k8sv1.DevelopmentDoguMap(*ldapDevelopmentDoguMap)
//
// 		managerWithMocks.resourceDoguFetcher.On("FetchWithResource", ctx, ldapCr).Return(ldapDogu, &developmentDoguMap, nil)
// 		managerWithMocks.imageRegistryMock.On("PullImageConfig", mock.Anything, mock.Anything).Return(imageConfig, nil)
// 		managerWithMocks.doguRegistratorMock.On("RegisterNewDogu", mock.Anything, mock.Anything, mock.Anything).Return(nil)
// 		managerWithMocks.dependencyValidatorMock.On("ValidateDependencies", ctx, mock.Anything).Return(nil)
// 		managerWithMocks.doguSecretHandlerMock.On("WriteDoguSecretsToRegistry", mock.Anything, mock.Anything).Return(nil)
// 		managerWithMocks.serviceAccountCreatorMock.On("CreateAll", mock.Anything, mock.Anything, mock.Anything).Return(nil)
// 		yamlResult := make(map[string]string, 0)
// 		managerWithMocks.fileExtractorMock.On("ExtractK8sResourcesFromContainer", mock.Anything, mock.Anything, mock.Anything).Return(yamlResult, nil)
// 		_ = managerWithMocks.installManager.clientset.Create(ctx, ldapCr)
// 		_ = managerWithMocks.installManager.clientset.Create(ctx, ldapDevelopmentDoguMap)
//
// 		managerWithMocks.recorder.On("Event", mock.Anything, corev1.EventTypeNormal, InstallEventReason, "Checking dependencies...").
// 			On("Event", mock.Anything, corev1.EventTypeNormal, InstallEventReason, "Registering in the local dogu registry...").
// 			On("Event", mock.Anything, corev1.EventTypeNormal, InstallEventReason, "Creating required service accounts...").
// 			On("Eventf", mock.Anything, corev1.EventTypeNormal, InstallEventReason, "Pulling dogu image %s...", "registry.cloudogu.com/official/ldap:2.4.48-4").
// 			On("Eventf", mock.Anything, corev1.EventTypeNormal, InstallEventReason, "Starting execPod...").
// 			On("Event", mock.Anything, corev1.EventTypeNormal, InstallEventReason, "Creating kubernetes resources...")
// 		managerWithMocks.resourceUpserter.On("UpsertDoguDeployment", ctx, ldapCr, ldapDogu, mock.Anything).Once().Return(nil, nil)
// 		managerWithMocks.resourceUpserter.On("UpsertDoguService", ctx, ldapCr, imageConfig).Once().Return(nil, nil)
// 		managerWithMocks.resourceUpserter.On("UpsertDoguExposedServices", ctx, ldapCr, ldapDogu).Once().Return(nil, nil)
// 		managerWithMocks.resourceUpserter.On("UpsertDoguPVCs", ctx, ldapCr, ldapDogu).Once().Return(nil, nil)
//
// 		execPod := mocks.NewExecPod(t)
// 		execPod.On("Create", testCtx).Return(nil)
// 		execPod.On("Delete", testCtx).Return(nil)
// 		managerWithMocks.execPodFactory.On("NewExecPod", internal.VolumeModeInstall, ldapCr, ldapDogu, mock.Anything).Return(execPod, nil)
//
// 		// when
// 		err := managerWithMocks.installManager.Install(ctx, ldapCr)
//
// 		// then
// 		require.NoError(t, err)
// 		managerWithMocks.AssertMocks(t)
//
// 		actualDevelopmentDoguMap := new(corev1.ConfigMap)
// 		err = managerWithMocks.installManager.clientset.Get(ctx, ldapCr.GetDevelopmentDoguMapKey(), actualDevelopmentDoguMap)
// 		require.True(t, apierrors.IsNotFound(err))
//
// 	})
//
// 	t.Run("failed to validate dependencies", func(t *testing.T) {
// 		// given
// 		managerWithMocks := getDoguInstallManagerWithMocks(t, getTestScheme())
// 		ldapCr, ldapDogu, _, _ := getDoguInstallManagerTestData(t)
// 		managerWithMocks.resourceDoguFetcher.On("FetchWithResource", ctx, ldapCr).Return(ldapDogu, nil, nil)
// 		managerWithMocks.dependencyValidatorMock.On("ValidateDependencies", ctx, mock.Anything).Return(assert.AnError)
// 		_ = managerWithMocks.installManager.clientset.Create(ctx, ldapCr)
//
// 		managerWithMocks.recorder.On("Event", mock.Anything, corev1.EventTypeNormal, InstallEventReason, "Checking dependencies...")
//
// 		// when
// 		err := managerWithMocks.installManager.Install(ctx, ldapCr)
//
// 		// then
// 		require.Error(t, err)
// 		assert.True(t, errors.Is(err, assert.AnError))
// 		managerWithMocks.AssertMocks(t)
// 	})
//
// 	t.Run("failed to register dogu", func(t *testing.T) {
// 		// given
// 		managerWithMocks := getDoguInstallManagerWithMocks(t, getTestScheme())
// 		ldapCr, _, _, _ := getDoguInstallManagerTestData(t)
// 		ldapCr, ldapDogu, _, _ := getDoguInstallManagerTestData(t)
// 		managerWithMocks.resourceDoguFetcher.On("FetchWithResource", ctx, ldapCr).Return(ldapDogu, nil, nil)
// 		managerWithMocks.doguRegistratorMock.On("RegisterNewDogu", mock.Anything, mock.Anything, mock.Anything).Return(assert.AnError)
// 		managerWithMocks.dependencyValidatorMock.On("ValidateDependencies", ctx, mock.Anything).Return(nil)
// 		_ = managerWithMocks.installManager.clientset.Create(ctx, ldapCr)
//
// 		managerWithMocks.recorder.On("Event", mock.Anything, corev1.EventTypeNormal, InstallEventReason, "Checking dependencies...")
// 		managerWithMocks.recorder.On("Event", mock.Anything, corev1.EventTypeNormal, InstallEventReason, "Registering in the local dogu registry...")
//
// 		// when
// 		err := managerWithMocks.installManager.Install(ctx, ldapCr)
//
// 		// then
// 		require.Error(t, err)
// 		assert.ErrorIs(t, err, assert.AnError)
// 		managerWithMocks.AssertMocks(t)
// 	})
//
// 	t.Run("failed to handle dogu secrets from setup", func(t *testing.T) {
// 		// given
// 		managerWithMocks := getDoguInstallManagerWithMocks(t, getTestScheme())
// 		ldapCr, _, _, _ := getDoguInstallManagerTestData(t)
// 		ldapCr, ldapDogu, _, _ := getDoguInstallManagerTestData(t)
// 		managerWithMocks.resourceDoguFetcher.On("FetchWithResource", ctx, ldapCr).Return(ldapDogu, nil, nil)
// 		managerWithMocks.doguRegistratorMock.On("RegisterNewDogu", mock.Anything, mock.Anything, mock.Anything).Return(nil)
// 		managerWithMocks.dependencyValidatorMock.On("ValidateDependencies", ctx, mock.Anything).Return(nil)
// 		managerWithMocks.doguSecretHandlerMock.On("WriteDoguSecretsToRegistry", mock.Anything, mock.Anything).Return(assert.AnError)
// 		_ = managerWithMocks.installManager.clientset.Create(ctx, ldapCr)
//
// 		managerWithMocks.recorder.On("Event", mock.Anything, corev1.EventTypeNormal, InstallEventReason, "Checking dependencies...")
// 		managerWithMocks.recorder.On("Event", mock.Anything, corev1.EventTypeNormal, InstallEventReason, "Registering in the local dogu registry...")
//
// 		// when
// 		err := managerWithMocks.installManager.Install(ctx, ldapCr)
//
// 		// then
// 		require.Error(t, err)
// 		assert.ErrorIs(t, err, assert.AnError)
// 		assert.ErrorContains(t, err, "failed to write dogu secrets from setup")
// 		managerWithMocks.AssertMocks(t)
// 	})
//
// 	t.Run("failed to create service accounts", func(t *testing.T) {
// 		// given
// 		managerWithMocks := getDoguInstallManagerWithMocks(t, getTestScheme())
// 		ldapCr, _, _, _ := getDoguInstallManagerTestData(t)
// 		ldapCr, ldapDogu, _, _ := getDoguInstallManagerTestData(t)
// 		managerWithMocks.resourceDoguFetcher.On("FetchWithResource", ctx, ldapCr).Return(ldapDogu, nil, nil)
// 		managerWithMocks.doguRegistratorMock.On("RegisterNewDogu", mock.Anything, mock.Anything, mock.Anything).Return(nil)
// 		managerWithMocks.dependencyValidatorMock.On("ValidateDependencies", ctx, mock.Anything).Return(nil)
// 		managerWithMocks.doguSecretHandlerMock.On("WriteDoguSecretsToRegistry", mock.Anything, mock.Anything).Return(nil)
// 		managerWithMocks.serviceAccountCreatorMock.On("CreateAll", mock.Anything, mock.Anything, mock.Anything).Return(assert.AnError)
// 		_ = managerWithMocks.installManager.clientset.Create(ctx, ldapCr)
//
// 		managerWithMocks.recorder.On("Event", mock.Anything, corev1.EventTypeNormal, InstallEventReason, "Checking dependencies...")
// 		managerWithMocks.recorder.On("Event", mock.Anything, corev1.EventTypeNormal, InstallEventReason, "Registering in the local dogu registry...")
// 		managerWithMocks.recorder.On("Event", mock.Anything, corev1.EventTypeNormal, InstallEventReason, "Creating required service accounts...")
//
// 		// when
// 		err := managerWithMocks.installManager.Install(ctx, ldapCr)
//
// 		// then
// 		require.Error(t, err)
// 		assert.ErrorIs(t, err, assert.AnError)
// 		assert.ErrorContains(t, err, "failed to create service accounts")
// 		managerWithMocks.AssertMocks(t)
// 	})
//
// 	t.Run("dogu resource not found", func(t *testing.T) {
// 		// given
// 		managerWithMocks := getDoguInstallManagerWithMocks(t, getTestScheme())
// 		ldapCr, _, _, _ := getDoguInstallManagerTestData(t)
//
// 		// when
// 		err := managerWithMocks.installManager.Install(ctx, ldapCr)
//
// 		// then
// 		require.Error(t, err)
// 		assert.ErrorContains(t, err, "not found")
// 		managerWithMocks.AssertMocks(t)
// 	})
//
// 	t.Run("error get dogu", func(t *testing.T) {
// 		// given
// 		managerWithMocks := getDoguInstallManagerWithMocks(t, getTestScheme())
// 		ldapCr, _, _, _ := getDoguInstallManagerTestData(t)
// 		managerWithMocks.resourceDoguFetcher.On("FetchWithResource", ctx, ldapCr).Return(nil, nil, assert.AnError)
//
// 		_ = managerWithMocks.installManager.clientset.Create(ctx, ldapCr)
//
// 		// when
// 		err := managerWithMocks.installManager.Install(ctx, ldapCr)
//
// 		// then
// 		require.Error(t, err)
// 		assert.ErrorIs(t, err, assert.AnError)
// 		managerWithMocks.AssertMocks(t)
// 	})
//
// 	t.Run("error on pull image", func(t *testing.T) {
// 		// given
// 		managerWithMocks := getDoguInstallManagerWithMocks(t, getTestScheme())
// 		ldapCr, _, _, _ := getDoguInstallManagerTestData(t)
// 		ldapCr, ldapDogu, _, _ := getDoguInstallManagerTestData(t)
// 		managerWithMocks.resourceDoguFetcher.On("FetchWithResource", ctx, ldapCr).Return(ldapDogu, nil, nil)
// 		managerWithMocks.imageRegistryMock.On("PullImageConfig", mock.Anything, mock.Anything).Return(nil, assert.AnError)
// 		managerWithMocks.doguRegistratorMock.On("RegisterNewDogu", mock.Anything, mock.Anything, mock.Anything).Return(nil)
// 		managerWithMocks.dependencyValidatorMock.On("ValidateDependencies", ctx, mock.Anything).Return(nil)
// 		managerWithMocks.doguSecretHandlerMock.On("WriteDoguSecretsToRegistry", mock.Anything, mock.Anything).Return(nil)
// 		managerWithMocks.serviceAccountCreatorMock.On("CreateAll", mock.Anything, mock.Anything, mock.Anything).Return(nil)
// 		_ = managerWithMocks.installManager.clientset.Create(ctx, ldapCr)
//
// 		managerWithMocks.recorder.On("Event", mock.Anything, corev1.EventTypeNormal, InstallEventReason, "Checking dependencies...")
// 		managerWithMocks.recorder.On("Event", mock.Anything, corev1.EventTypeNormal, InstallEventReason, "Registering in the local dogu registry...")
// 		managerWithMocks.recorder.On("Event", mock.Anything, corev1.EventTypeNormal, InstallEventReason, "Creating required service accounts...")
// 		managerWithMocks.recorder.On("Eventf", mock.Anything, corev1.EventTypeNormal, InstallEventReason, "Pulling dogu image %s...", "registry.cloudogu.com/official/ldap:2.4.48-4")
//
// 		// when
// 		err := managerWithMocks.installManager.Install(ctx, ldapCr)
//
// 		// then
// 		require.Error(t, err)
// 		assert.ErrorIs(t, err, assert.AnError)
// 		managerWithMocks.AssertMocks(t)
// 	})
//
// 	t.Run("error on upsert", func(t *testing.T) {
// 		t.Run("succeeds", func(t *testing.T) {
// 			// given
// 			managerWithMocks := getDoguInstallManagerWithMocks(t, getTestScheme())
// 			ldapCr, ldapDogu, _, imageConfig := getDoguInstallManagerTestData(t)
// 			managerWithMocks.resourceDoguFetcher.On("FetchWithResource", ctx, ldapCr).Return(ldapDogu, nil, nil)
// 			managerWithMocks.imageRegistryMock.On("PullImageConfig", mock.Anything, mock.Anything).Return(imageConfig, nil)
// 			managerWithMocks.doguRegistratorMock.On("RegisterNewDogu", mock.Anything, mock.Anything, mock.Anything).Return(nil)
// 			managerWithMocks.dependencyValidatorMock.On("ValidateDependencies", ctx, mock.Anything).Return(nil)
// 			managerWithMocks.doguSecretHandlerMock.On("WriteDoguSecretsToRegistry", mock.Anything, mock.Anything).Return(nil)
// 			managerWithMocks.serviceAccountCreatorMock.On("CreateAll", mock.Anything, mock.Anything, mock.Anything).Return(nil)
// 			yamlResult := make(map[string]string, 0)
// 			managerWithMocks.fileExtractorMock.On("ExtractK8sResourcesFromContainer", mock.Anything, mock.Anything, mock.Anything).Return(yamlResult, nil)
// 			ldapCr.ResourceVersion = ""
// 			_ = managerWithMocks.installManager.clientset.Create(ctx, ldapCr)
//
// 			managerWithMocks.resourceUpserter.On("UpsertDoguDeployment", ctx, ldapCr, ldapDogu, mock.Anything).Once().Return(nil, nil)
// 			managerWithMocks.resourceUpserter.On("UpsertDoguService", ctx, ldapCr, imageConfig).Once().Return(nil, nil)
// 			managerWithMocks.resourceUpserter.On("UpsertDoguExposedServices", ctx, ldapCr, ldapDogu).Once().Return(nil, nil)
// 			managerWithMocks.resourceUpserter.On("UpsertDoguPVCs", ctx, ldapCr, ldapDogu).Once().Return(nil, nil)
//
// 			managerWithMocks.recorder.On("Event", mock.Anything, corev1.EventTypeNormal, InstallEventReason, "Checking dependencies...").
// 				On("Event", mock.Anything, corev1.EventTypeNormal, InstallEventReason, "Registering in the local dogu registry...").
// 				On("Event", mock.Anything, corev1.EventTypeNormal, InstallEventReason, "Creating required service accounts...").
// 				On("Eventf", mock.Anything, corev1.EventTypeNormal, InstallEventReason, "Pulling dogu image %s...", "registry.cloudogu.com/official/ldap:2.4.48-4").
// 				On("Eventf", mock.Anything, corev1.EventTypeNormal, InstallEventReason, "Starting execPod...").
// 				On("Event", mock.Anything, corev1.EventTypeNormal, InstallEventReason, "Creating kubernetes resources...")
// 			execPod := mocks.NewExecPod(t)
// 			execPod.On("Create", testCtx).Return(nil)
// 			execPod.On("Delete", testCtx).Return(nil)
// 			managerWithMocks.execPodFactory.On("NewExecPod", internal.VolumeModeInstall, ldapCr, ldapDogu, mock.Anything).Return(execPod, nil)
//
// 			// when
// 			err := managerWithMocks.installManager.Install(ctx, ldapCr)
//
// 			// then
// 			require.NoError(t, err)
// 			managerWithMocks.AssertMocks(t)
// 		})
// 		t.Run("fails when upserting deployment", func(t *testing.T) {
// 			// given
// 			managerWithMocks := getDoguInstallManagerWithMocks(t, getTestScheme())
// 			ldapCr, ldapDogu, _, imageConfig := getDoguInstallManagerTestData(t)
// 			managerWithMocks.resourceDoguFetcher.On("FetchWithResource", ctx, ldapCr).Return(ldapDogu, nil, nil)
// 			managerWithMocks.imageRegistryMock.On("PullImageConfig", mock.Anything, mock.Anything).Return(imageConfig, nil)
// 			managerWithMocks.doguRegistratorMock.On("RegisterNewDogu", mock.Anything, mock.Anything, mock.Anything).Return(nil)
// 			managerWithMocks.dependencyValidatorMock.On("ValidateDependencies", ctx, mock.Anything).Return(nil)
// 			managerWithMocks.doguSecretHandlerMock.On("WriteDoguSecretsToRegistry", mock.Anything, mock.Anything).Return(nil)
// 			managerWithMocks.serviceAccountCreatorMock.On("CreateAll", mock.Anything, mock.Anything, mock.Anything).Return(nil)
// 			yamlResult := make(map[string]string, 0)
// 			managerWithMocks.fileExtractorMock.On("ExtractK8sResourcesFromContainer", mock.Anything, mock.Anything, mock.Anything).Return(yamlResult, nil)
// 			ldapCr.ResourceVersion = ""
// 			_ = managerWithMocks.installManager.clientset.Create(ctx, ldapCr)
//
// 			managerWithMocks.recorder.On("Event", mock.Anything, corev1.EventTypeNormal, InstallEventReason, "Checking dependencies...").
// 				On("Event", mock.Anything, corev1.EventTypeNormal, InstallEventReason, "Registering in the local dogu registry...").
// 				On("Event", mock.Anything, corev1.EventTypeNormal, InstallEventReason, "Creating required service accounts...").
// 				On("Eventf", mock.Anything, corev1.EventTypeNormal, InstallEventReason, "Pulling dogu image %s...", "registry.cloudogu.com/official/ldap:2.4.48-4").
// 				On("Eventf", mock.Anything, corev1.EventTypeNormal, InstallEventReason, "Starting execPod...").
// 				On("Event", mock.Anything, corev1.EventTypeNormal, InstallEventReason, "Creating kubernetes resources...")
// 			execPod := mocks.NewExecPod(t)
// 			execPod.On("Create", testCtx).Return(nil)
// 			execPod.On("Delete", testCtx).Return(nil)
// 			managerWithMocks.execPodFactory.On("NewExecPod", internal.VolumeModeInstall, ldapCr, ldapDogu, mock.Anything).Return(execPod, nil)
//
// 			managerWithMocks.resourceUpserter.On("UpsertDoguService", ctx, ldapCr, imageConfig).Once().Return(nil, nil)
// 			managerWithMocks.resourceUpserter.On("UpsertDoguExposedServices", ctx, ldapCr, ldapDogu).Once().Return(nil, nil)
// 			managerWithMocks.resourceUpserter.On("UpsertDoguDeployment", ctx, ldapCr, ldapDogu, mock.Anything).Once().Return(nil, assert.AnError)
//
// 			// when
// 			err := managerWithMocks.installManager.Install(ctx, ldapCr)
//
// 			// then
// 			require.Error(t, err)
// 			assert.ErrorContains(t, err, "failed to create dogu resources")
// 			assert.ErrorIs(t, err, assert.AnError)
// 			managerWithMocks.AssertMocks(t)
// 		})
// 		t.Run("fails when upserting service", func(t *testing.T) {
// 			// given
// 			managerWithMocks := getDoguInstallManagerWithMocks(t, getTestScheme())
// 			ldapCr, ldapDogu, _, imageConfig := getDoguInstallManagerTestData(t)
// 			managerWithMocks.resourceDoguFetcher.On("FetchWithResource", ctx, ldapCr).Return(ldapDogu, nil, nil)
// 			managerWithMocks.imageRegistryMock.On("PullImageConfig", mock.Anything, mock.Anything).Return(imageConfig, nil)
// 			managerWithMocks.doguRegistratorMock.On("RegisterNewDogu", mock.Anything, mock.Anything, mock.Anything).Return(nil)
// 			managerWithMocks.dependencyValidatorMock.On("ValidateDependencies", mock.Anything, ldapDogu).Return(nil)
// 			managerWithMocks.doguSecretHandlerMock.On("WriteDoguSecretsToRegistry", mock.Anything, mock.Anything).Return(nil)
// 			managerWithMocks.serviceAccountCreatorMock.On("CreateAll", mock.Anything, mock.Anything, mock.Anything).Return(nil)
// 			ldapCr.ResourceVersion = ""
// 			_ = managerWithMocks.installManager.clientset.Create(ctx, ldapCr)
//
// 			managerWithMocks.recorder.On("Event", mock.Anything, corev1.EventTypeNormal, InstallEventReason, "Checking dependencies...").
// 				On("Event", mock.Anything, corev1.EventTypeNormal, InstallEventReason, "Registering in the local dogu registry...").
// 				On("Event", mock.Anything, corev1.EventTypeNormal, InstallEventReason, "Creating required service accounts...").
// 				On("Eventf", mock.Anything, corev1.EventTypeNormal, InstallEventReason, "Pulling dogu image %s...", "registry.cloudogu.com/official/ldap:2.4.48-4").
// 				On("Event", mock.Anything, corev1.EventTypeNormal, InstallEventReason, "Creating kubernetes resources...")
// 			managerWithMocks.resourceUpserter.On("UpsertDoguService", ctx, ldapCr, imageConfig).Once().Return(nil, assert.AnError)
//
// 			// when
// 			err := managerWithMocks.installManager.Install(ctx, ldapCr)
//
// 			// then
// 			require.Error(t, err)
// 			assert.ErrorContains(t, err, "failed to create dogu resources")
// 			assert.ErrorIs(t, err, assert.AnError)
// 			managerWithMocks.AssertMocks(t)
// 		})
// 		t.Run("fails when upserting exposed services", func(t *testing.T) {
// 			// given
// 			managerWithMocks := getDoguInstallManagerWithMocks(t, getTestScheme())
// 			ldapCr, ldapDogu, _, imageConfig := getDoguInstallManagerTestData(t)
// 			managerWithMocks.resourceDoguFetcher.On("FetchWithResource", ctx, ldapCr).Return(ldapDogu, nil, nil)
// 			managerWithMocks.imageRegistryMock.On("PullImageConfig", mock.Anything, mock.Anything).Return(imageConfig, nil)
// 			managerWithMocks.doguRegistratorMock.On("RegisterNewDogu", mock.Anything, mock.Anything, mock.Anything).Return(nil)
// 			managerWithMocks.dependencyValidatorMock.On("ValidateDependencies", mock.Anything, ldapDogu).Return(nil)
// 			managerWithMocks.doguSecretHandlerMock.On("WriteDoguSecretsToRegistry", mock.Anything, mock.Anything).Return(nil)
// 			managerWithMocks.serviceAccountCreatorMock.On("CreateAll", mock.Anything, mock.Anything, mock.Anything).Return(nil)
// 			ldapCr.ResourceVersion = ""
// 			_ = managerWithMocks.installManager.clientset.Create(ctx, ldapCr)
//
// 			managerWithMocks.recorder.On("Event", mock.Anything, corev1.EventTypeNormal, InstallEventReason, "Checking dependencies...").
// 				On("Event", mock.Anything, corev1.EventTypeNormal, InstallEventReason, "Registering in the local dogu registry...").
// 				On("Event", mock.Anything, corev1.EventTypeNormal, InstallEventReason, "Creating required service accounts...").
// 				On("Eventf", mock.Anything, corev1.EventTypeNormal, InstallEventReason, "Pulling dogu image %s...", "registry.cloudogu.com/official/ldap:2.4.48-4").
// 				On("Event", mock.Anything, corev1.EventTypeNormal, InstallEventReason, "Creating kubernetes resources...")
// 			managerWithMocks.resourceUpserter.On("UpsertDoguService", ctx, ldapCr, imageConfig).Once().Return(nil, nil)
// 			managerWithMocks.resourceUpserter.On("UpsertDoguExposedServices", ctx, ldapCr, ldapDogu).Once().Return(nil, assert.AnError)
//
// 			// when
// 			err := managerWithMocks.installManager.Install(ctx, ldapCr)
//
// 			// then
// 			require.Error(t, err)
// 			assert.ErrorContains(t, err, "failed to create dogu resources")
// 			assert.ErrorIs(t, err, assert.AnError)
// 			managerWithMocks.AssertMocks(t)
// 		})
// 		t.Run("fails when upserting pvcs", func(t *testing.T) {
// 			// given
// 			managerWithMocks := getDoguInstallManagerWithMocks(t, getTestScheme())
// 			ldapCr, ldapDogu, _, imageConfig := getDoguInstallManagerTestData(t)
// 			managerWithMocks.resourceDoguFetcher.On("FetchWithResource", ctx, ldapCr).Return(ldapDogu, nil, nil)
// 			managerWithMocks.imageRegistryMock.On("PullImageConfig", mock.Anything, mock.Anything).Return(imageConfig, nil)
// 			managerWithMocks.doguRegistratorMock.On("RegisterNewDogu", mock.Anything, mock.Anything, mock.Anything).Return(nil)
// 			managerWithMocks.dependencyValidatorMock.On("ValidateDependencies", mock.Anything, ldapDogu).Return(nil)
// 			managerWithMocks.doguSecretHandlerMock.On("WriteDoguSecretsToRegistry", mock.Anything, mock.Anything).Return(nil)
// 			managerWithMocks.serviceAccountCreatorMock.On("CreateAll", mock.Anything, mock.Anything, mock.Anything).Return(nil)
// 			yamlResult := make(map[string]string, 0)
// 			managerWithMocks.fileExtractorMock.On("ExtractK8sResourcesFromContainer", mock.Anything, mock.Anything, mock.Anything).Return(yamlResult, nil)
// 			ldapCr.ResourceVersion = ""
// 			_ = managerWithMocks.installManager.clientset.Create(ctx, ldapCr)
//
// 			managerWithMocks.recorder.On("Event", mock.Anything, corev1.EventTypeNormal, InstallEventReason, "Checking dependencies...").
// 				On("Event", mock.Anything, corev1.EventTypeNormal, InstallEventReason, "Registering in the local dogu registry...").
// 				On("Event", mock.Anything, corev1.EventTypeNormal, InstallEventReason, "Creating required service accounts...").
// 				On("Eventf", mock.Anything, corev1.EventTypeNormal, InstallEventReason, "Pulling dogu image %s...", "registry.cloudogu.com/official/ldap:2.4.48-4").
// 				On("Eventf", mock.Anything, corev1.EventTypeNormal, InstallEventReason, "Starting execPod...").
// 				On("Event", mock.Anything, corev1.EventTypeNormal, InstallEventReason, "Creating kubernetes resources...")
// 			execPod := mocks.NewExecPod(t)
// 			execPod.On("Create", testCtx).Return(nil)
// 			execPod.On("Delete", testCtx).Return(nil)
// 			managerWithMocks.execPodFactory.On("NewExecPod", internal.VolumeModeInstall, ldapCr, ldapDogu, mock.Anything).Return(execPod, nil)
//
// 			managerWithMocks.resourceUpserter.On("UpsertDoguDeployment", ctx, ldapCr, ldapDogu, mock.Anything).Once().Return(nil, nil)
// 			managerWithMocks.resourceUpserter.On("UpsertDoguService", ctx, ldapCr, imageConfig).Once().Return(nil, nil)
// 			managerWithMocks.resourceUpserter.On("UpsertDoguExposedServices", ctx, ldapCr, ldapDogu).Once().Return(nil, nil)
// 			managerWithMocks.resourceUpserter.On("UpsertDoguPVCs", ctx, ldapCr, ldapDogu).Once().Return(nil, assert.AnError)
//
// 			// when
// 			err := managerWithMocks.installManager.Install(ctx, ldapCr)
//
// 			// then
// 			require.Error(t, err)
// 			assert.ErrorContains(t, err, "failed to create dogu resources")
// 			assert.ErrorIs(t, err, assert.AnError)
// 			managerWithMocks.AssertMocks(t)
// 		})
// 	})
// }
