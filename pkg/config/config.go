package config

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/cloudogu/cesapp-lib/core"

	"gopkg.in/yaml.v3"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"
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

// HelmRepositoryData contains all necessary data for the helm repository.
type HelmRepositoryData struct {
	// Endpoint contains the Helm registry endpoint URL.
	Endpoint string `json:"endpoint" yaml:"endpoint"`
	// PlainHttp indicates that the repository endpoint should be accessed using plain http
	PlainHttp bool `json:"plainHttp,omitempty" yaml:"plainHttp,omitempty"`
}

// GetOciEndpoint returns the configured endpoint of the HelmRepositoryData with the OCI-protocol
func (hrd *HelmRepositoryData) GetOciEndpoint() (string, error) {
	split := strings.Split(hrd.Endpoint, "://")
	if len(split) == 1 && split[0] != "" {
		return fmt.Sprintf("oci://%s", split[0]), nil
	}
	if len(split) == 2 && split[1] != "" {
		return fmt.Sprintf("oci://%s", split[1]), nil
	}

	return "", fmt.Errorf("error creating oci-endpoint from '%s': wrong format", hrd.Endpoint)
}

func (hrd *HelmRepositoryData) IsPlainHttp() bool {
	return hrd.PlainHttp
}

// OperatorConfig contains all configurable values for the dogu operator.
type OperatorConfig struct {
	// Namespace specifies the namespace that the operator is deployed to.
	Namespace string `json:"namespace"`
	// Version contains the current version of the operator
	Version *core.Version `json:"version"`
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

	parsedVersion, err := core.ParseVersion(version)
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
		Version:   &parsedVersion,
	}, nil
}

// GetHelmRepositoryData reads the repository data either from file or from a secret in the cluster.
func GetHelmRepositoryData(configMapClient corev1.ConfigMapInterface) (*HelmRepositoryData, error) {
	runtime, err := getEnvVar(runtimeEnvironmentVariable)
	if err != nil {
		log.Info("Runtime env var not found.")
	}

	if runtime == runtimeLocal {
		return getHelmRepositoryDataFromFile()
	}

	return getHelmRepositoryFromConfigMap(configMapClient)
}

func getHelmRepositoryFromConfigMap(configMapClient corev1.ConfigMapInterface) (*HelmRepositoryData, error) {
	configMap, err := configMapClient.Get(context.TODO(), helmRepositoryConfigMapName, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		return nil, fmt.Errorf("helm repository configMap %s not found: %w", helmRepositoryConfigMapName, err)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get helm repository configMap %s: %w", helmRepositoryConfigMapName, err)
	}

	plainHttp := false
	const configMapPlainHttp = "plainHttp"
	if plainHttpStr, exists := configMap.Data[configMapPlainHttp]; exists {
		plainHttp, err = strconv.ParseBool(plainHttpStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse field %s from configMap %s", configMapPlainHttp, helmRepositoryConfigMapName)
		}
	}

	return &HelmRepositoryData{
		Endpoint:  configMap.Data["endpoint"],
		PlainHttp: plainHttp,
	}, nil
}

func getHelmRepositoryDataFromFile() (*HelmRepositoryData, error) {
	data := &HelmRepositoryData{}
	if _, err := os.Stat(devHelmRepoDataPath); os.IsNotExist(err) {
		return data, fmt.Errorf("could not find configuration at %s", devHelmRepoDataPath)
	}

	fileData, err := os.ReadFile(devHelmRepoDataPath)
	if err != nil {
		return data, fmt.Errorf("failed to read configuration %s: %w", devHelmRepoDataPath, err)
	}

	err = yaml.Unmarshal(fileData, data)
	if err != nil {
		return data, fmt.Errorf("failed to unmarshal configuration %s: %w", devHelmRepoDataPath, err)
	}

	return data, nil
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
