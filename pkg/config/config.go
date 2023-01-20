package config

import (
	"fmt"
	"github.com/cloudogu/cesapp-lib/core"
	"os"
	ctrl "sigs.k8s.io/controller-runtime"
	"strings"
)

const (
	StageDevelopment         = "development"
	StageProduction          = "production"
	StageEnvironmentVariable = "STAGE"
)

var Stage = StageProduction

var (
	envVarNamespace             = "NAMESPACE"
	envVarDoguRegistryEndpoint  = "DOGU_REGISTRY_ENDPOINT"
	envVarDoguRegistryUsername  = "DOGU_REGISTRY_USERNAME"
	envVarDoguRegistryPassword  = "DOGU_REGISTRY_PASSWORD"
	envVarDoguRegistryURLSchema = "DOGU_REGISTRY_URLSCHEMA"
	log                         = ctrl.Log.WithName("config")
)

type DoguRegistryData struct {
	Endpoint  string `json:"endpoint"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	URLSchema string `json:"urlschema"`
}

// OperatorConfig contains all configurable values for the dogu operator.
type OperatorConfig struct {
	// Namespace specifies the namespace that the operator is deployed to.
	Namespace string `json:"namespace"`
	// Version contains the current version of the operator
	Version *core.Version `json:"version"`
	// DoguRegistry contains all necessary data for the dogu registry.
	DoguRegistry DoguRegistryData `json:"dogu_registry"`
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
	// doguRegistryData, err := readDoguRegistryData()
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to read dogu registry data: %w", err)
	// }
	// log.Info(fmt.Sprintf("Found stored dogu registry data! Using dogu registry %s", doguRegistryData.Endpoint))

	return &OperatorConfig{
		Namespace:    namespace,
		Version:      &parsedVersion,
		DoguRegistry: DoguRegistryData{},
	}, nil
}

func readDoguRegistryData() (DoguRegistryData, error) {
	endpoint, err := getEnvVar(envVarDoguRegistryEndpoint)
	if err != nil {
		return DoguRegistryData{}, fmt.Errorf("failed to get env var [%s]: %w", envVarDoguRegistryEndpoint, err)
	}
	// remove tailing slash
	endpoint = strings.TrimSuffix(endpoint, "/")

	username, err := getEnvVar(envVarDoguRegistryUsername)
	if err != nil {
		return DoguRegistryData{}, fmt.Errorf("failed to get env var [%s]: %w", envVarDoguRegistryUsername, err)
	}

	password, err := getEnvVar(envVarDoguRegistryPassword)
	if err != nil {
		return DoguRegistryData{}, fmt.Errorf("failed to get env var [%s]: %w", envVarDoguRegistryPassword, err)
	}

	urlschema, err := getEnvVar(envVarDoguRegistryURLSchema)
	if err != nil {
		log.Info(envVarDoguRegistryURLSchema + " not set, using default")
		urlschema = "default"
	}

	return DoguRegistryData{
		Endpoint:  endpoint,
		Username:  username,
		Password:  password,
		URLSchema: urlschema,
	}, nil
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
