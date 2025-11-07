package config

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/Masterminds/semver/v3"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/yaml"
)

const (
	// StageDevelopment is the name for the development-stage
	StageDevelopment = "development"
	// StageProduction is the name for the production-stage
	StageProduction = "production"
	// StageEnvironmentVariable is the name of the environment-variable containing the configured stage
	StageEnvironmentVariable = "STAGE"
	// runtimeEnvironmentVariable is the name of the environment-variable containing the configured runtime
	runtimeEnvironmentVariable = "RUNTIME"
	// runtimeLocal is the name for the local-runtime on a developer-machine
	runtimeLocal = "local"
	// RequeueTimeInNanosecondsEnvironmentVariable is the name of the environment variable containing the configured requeueTime
	RequeueTimeInNanosecondsEnvironmentVariable = "REQUEUE_TIME_IN_NANOSECONDS"
	// helmRepositoryConfigMapName is the name
	helmRepositoryConfigMapName = "component-operator-helm-repository"
)

const defaultRequeueTime = time.Second * 3

var (
	Stage               = StageProduction
	devHelmRepoDataPath = "k8s/helm-repository.yaml"
)

var (
	envVarNamespace = "NAMESPACE"

	envHelmClientTimeoutMins      = "HELM_CLIENT_TIMEOUT_MINS"
	defaultHelmClientTimeoutMins  = time.Duration(15) * time.Minute
	envHealthSyncIntervalMins     = "HEALTH_SYNC_INTERVAL_MINS"
	defaultHealthSyncIntervalMins = time.Duration(2) * time.Minute

	log = ctrl.Log.WithName("config")
)

const (
	configMapSchema      = "schema"
	configMapPlainHttp   = "plainHttp"
	configMapInsecureTls = "insecureTls"
)

type EndpointSchema string

const EndpointSchemaOCI EndpointSchema = "oci"

type configMapInterface interface {
	corev1.ConfigMapInterface
}

// HelmRepositoryData contains all necessary data for the helm repository.
type HelmRepositoryData struct {
	// Endpoint contains the Helm registry endpoint URL.
	Endpoint string `json:"endpoint" yaml:"endpoint"`
	// Schema describes the way how clients communicate with the Helm registry endpoint.
	Schema EndpointSchema `json:"schema" yaml:"schema"`
	// PlainHttp indicates that the repository endpoint should be accessed using plain http
	PlainHttp bool `json:"plainHttp,omitempty" yaml:"plainHttp,omitempty"`
	// InsecureTls allows invalid or selfsigned certificates to be used. This option may be overridden by PlainHttp which forces HTTP traffic.
	InsecureTLS bool `json:"insecureTls" yaml:"insecureTls"`
}

// URL returns the full URL Helm repository endpoint including schema.
func (hrd *HelmRepositoryData) URL() string {
	input := []string{string(hrd.Schema), hrd.Endpoint}

	return strings.Join(input, "://")
}

func (hrd *HelmRepositoryData) validate() error {
	if hrd.Endpoint == "" {
		return fmt.Errorf("endpoint URL must not be empty")
	}

	if strings.Contains(hrd.Endpoint, "://") {
		return fmt.Errorf("endpoint URL '%s' solely consist of the endpoint without schema or ://", hrd.Endpoint)
	}

	if hrd.Schema != EndpointSchemaOCI {
		return fmt.Errorf("endpoint uses an unsupported schema '%s': valid schemas are: oci", hrd.Schema)
	}

	return nil
}

// OperatorConfig contains all configurable values for the component operator.
type OperatorConfig struct {
	// Namespace specifies the namespace that the operator is deployed to.
	Namespace string `json:"namespace"`
	// Version contains the current version of the operator
	Version *semver.Version `json:"version"`
	// HelmRepositoryData contains all necessary data for the helm repository.
	HelmRepositoryData     *HelmRepositoryData `json:"helm_repository"`
	HelmClientTimeoutMins  time.Duration
	HealthSyncIntervalMins time.Duration
	RequeueTime            time.Duration
}

// NewOperatorConfig creates a new operator config by reading values from the environment variables
func NewOperatorConfig(version string) (*OperatorConfig, error) {
	stage, err := getEnvVar(StageEnvironmentVariable)
	if err != nil {
		log.Error(err, "Error reading stage environment variable. Use Stage production")
	}
	Stage = stage

	if Stage == StageDevelopment {
		log.Info("Starting in development mode! This is not recommended for production!")
	}

	parsedVersion, err := semver.NewVersion(version)
	if err != nil {
		return nil, fmt.Errorf("failed to parse version: %w", err)
	}
	log.Info(fmt.Sprintf("Version: [%s]", version))

	namespace, err := readNamespace()
	if err != nil {
		return nil, fmt.Errorf("failed to read namespace: %w", err)
	}
	log.Info(fmt.Sprintf("Deploying the k8s component operator in namespace %s", namespace))

	requeueTime, err := readReconcilerRequeueTime()
	if err != nil {
		log.Error(err, fmt.Sprintf("failed to read requeue time. Using default requeue time %s", defaultRequeueTime))
	}

	return &OperatorConfig{
		Namespace:              namespace,
		Version:                parsedVersion,
		HelmClientTimeoutMins:  readMinuteDurationEnv(envHelmClientTimeoutMins, defaultHelmClientTimeoutMins),
		HealthSyncIntervalMins: readMinuteDurationEnv(envHealthSyncIntervalMins, defaultHealthSyncIntervalMins),
		RequeueTime:            requeueTime,
	}, nil
}

// GetHelmRepositoryData reads the repository data either from file or from a secret in the cluster.
func GetHelmRepositoryData(ctx context.Context, configMapClient configMapInterface) (*HelmRepositoryData, error) {
	runtime, err := getEnvVar(runtimeEnvironmentVariable)
	if err != nil {
		log.Info("Runtime env var not found.")
	}

	if runtime == runtimeLocal {
		return NewHelmRepoDataFromFile(devHelmRepoDataPath)
	}

	return NewHelmRepoDataFromCluster(ctx, configMapClient)
}

// NewHelmRepoDataFromCluster reads the repo data ConfigMap, validates and returns it.
func NewHelmRepoDataFromCluster(ctx context.Context, configMapClient configMapInterface) (*HelmRepositoryData, error) {
	configMap, err := configMapClient.Get(ctx, helmRepositoryConfigMapName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get helm repository configMap %s: %w", helmRepositoryConfigMapName, err)
	}

	plainHttp := false
	if plainHttpStr, exists := configMap.Data[configMapPlainHttp]; exists {
		plainHttp, err = strconv.ParseBool(plainHttpStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse field %s from configMap %s", configMapPlainHttp, helmRepositoryConfigMapName)
		}
	}
	insecureTls := false
	if insecureTlsStr, exists := configMap.Data[configMapInsecureTls]; exists {
		insecureTls, err = strconv.ParseBool(insecureTlsStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse field %s from configMap %s", configMapInsecureTls, helmRepositoryConfigMapName)
		}
	}

	schema := configMap.Data[configMapSchema]
	repoData := &HelmRepositoryData{
		Endpoint:    configMap.Data["endpoint"],
		Schema:      EndpointSchema(schema),
		PlainHttp:   plainHttp,
		InsecureTLS: insecureTls,
	}

	err = repoData.validate()
	if err != nil {
		return nil, fmt.Errorf("config map '%s' failed validation: %w", helmRepositoryConfigMapName, err)
	}

	return repoData, nil
}

func NewHelmRepoDataFromFile(filepath string) (*HelmRepositoryData, error) {
	fileBytes, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read configuration %s: %w", filepath, err)
	}

	repoData := &HelmRepositoryData{}
	err = yaml.Unmarshal(fileBytes, repoData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal configuration %s: %w", filepath, err)
	}

	err = repoData.validate()
	if err != nil {
		return nil, fmt.Errorf("helm repository data from file '%s' failed validation: %w", filepath, err)
	}

	return repoData, nil
}

func readNamespace() (string, error) {
	namespace, err := getEnvVar(envVarNamespace)
	if err != nil {
		return "", err
	}

	return namespace, nil
}

func getEnvVar(name string) (string, error) {
	ns, found := os.LookupEnv(name)
	if !found {
		return "", fmt.Errorf("environment variable %s must be set", name)
	}
	return ns, nil
}

func readMinuteDurationEnv(env string, defaultValue time.Duration) time.Duration {
	valueString, err := getEnvVar(env)
	if err != nil {
		logrus.Warningf("failed to read %s environment variable, using default value", env)
		return defaultValue
	}

	valueParsed, err := strconv.Atoi(valueString)
	if err != nil {
		logrus.Warningf("failed to parse %s environment variable, using default value", env)
		return defaultValue
	}

	if valueParsed <= 0 {
		logrus.Warningf("parsed value (%d) of %s is smaller than 0, using default value", valueParsed, env)
		return defaultValue
	}

	return time.Duration(valueParsed) * time.Minute
}

func readReconcilerRequeueTime() (time.Duration, error) {
	requeueTimeString, err := getEnvVar(RequeueTimeInNanosecondsEnvironmentVariable)
	if err != nil {
		return defaultRequeueTime, newEnvVarError(envVarNamespace, err)
	}
	requeueTime, err := strconv.ParseFloat(requeueTimeString, 64)
	if err != nil {
		return defaultRequeueTime, err
	}
	return time.Duration(requeueTime), nil
}

func newEnvVarError(envVar string, err error) error {
	return fmt.Errorf("failed to get env var [%s]: %w", envVar, err)
}
