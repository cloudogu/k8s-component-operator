package config

import (
	"context"
	"fmt"
	"github.com/Masterminds/semver/v3"
	"gopkg.in/yaml.v3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"os"
	ctrl "sigs.k8s.io/controller-runtime"
	"strconv"
	"strings"
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
	// helmRepositoryConfigMapName is the name
	helmRepositoryConfigMapName = "component-operator-helm-repository"
)

var (
	Stage               = StageProduction
	devHelmRepoDataPath = "k8s/helm-repository.yaml"
)

var (
	envVarNamespace = "NAMESPACE"
	log             = ctrl.Log.WithName("config")
)

const (
	configMapSchema    = "schema"
	configMapPlainHttp = "plainHttp"
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

	if hrd.Schema != EndpointSchemaOCI {
		return fmt.Errorf("endpoint uses an unsupported schema '%s': valid schemas are: oci", hrd.Schema)
	}

	return nil
}

// OperatorConfig contains all configurable values for the dogu operator.
type OperatorConfig struct {
	// Namespace specifies the namespace that the operator is deployed to.
	Namespace string `json:"namespace"`
	// Version contains the current version of the operator
	Version *semver.Version `json:"version"`
	// HelmRepositoryData contains all necessary data for the helm repository.
	HelmRepositoryData *HelmRepositoryData `json:"helm_repository"`
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
	log.Info(fmt.Sprintf("Deploying the k8s dogu operator in namespace %s", namespace))

	return &OperatorConfig{
		Namespace: namespace,
		Version:   parsedVersion,
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
	var schema string
	var schemaExists bool
	if schema, schemaExists = configMap.Data[configMapSchema]; schemaExists {
		if !schemaExists {
			return nil, fmt.Errorf("field %s does not exist in configMap %s", configMapSchema, helmRepositoryConfigMapName)
		}
	}

	repoData := &HelmRepositoryData{
		Endpoint:  configMap.Data["endpoint"],
		Schema:    EndpointSchema(schema),
		PlainHttp: plainHttp,
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
		return nil, fmt.Errorf("helm repository data from file '%s' failed validation: %w", fileBytes, err)
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
