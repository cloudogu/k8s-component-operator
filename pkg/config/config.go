package config

import (
	"context"
	"fmt"
	"github.com/cloudogu/cesapp-lib/core"
	"gopkg.in/yaml.v3"
	"k8s.io/apimachinery/pkg/api/errors"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"os"
	ctrl "sigs.k8s.io/controller-runtime"
	"strings"
)

const (
	StageDevelopment            = "development"
	StageProduction             = "production"
	StageEnvironmentVariable    = "STAGE"
	runtimeEnvironmentVariable  = "RUNTIME"
	runtimeLocal                = "local"
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

type HelmRepositoryData struct {
	Endpoint string `json:"endpoint"`
}

func (hrd *HelmRepositoryData) GetOciEndpoint() (string, error) {
	split := strings.Split(hrd.Endpoint, "://")
	if len(split) != 2 {
		return "", fmt.Errorf("error creating oci-endpoint from '%s': wrong format", hrd.Endpoint)
	}
	return fmt.Sprintf("oci://%s", split[1]), nil
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
func GetHelmRepositoryData(configMapClient v1.ConfigMapInterface) (*HelmRepositoryData, error) {
	runtime, err := getEnvVar(runtimeEnvironmentVariable)
	if err != nil {
		log.Info("Runtime env var not found.")
	}

	if runtime == runtimeLocal {
		return getHelmRepositoryDataFromFile()
	} else {
		return getHelmRepositoryFromConfigMap(configMapClient)
	}
}

func getHelmRepositoryFromConfigMap(configMapClient v1.ConfigMapInterface) (*HelmRepositoryData, error) {
	configMap, err := configMapClient.Get(context.TODO(), helmRepositoryConfigMapName, v12.GetOptions{})
	if errors.IsNotFound(err) {
		return nil, fmt.Errorf("helm repository configMap %s not found: %w", helmRepositoryConfigMapName, err)
	} else if err != nil {
		return nil, fmt.Errorf("failed to get helm repository configMap %s: %w", helmRepositoryConfigMapName, err)
	}

	return &HelmRepositoryData{
		Endpoint: configMap.Data["endpoint"],
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
