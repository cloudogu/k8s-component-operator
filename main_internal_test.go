package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"testing"

	v1 "github.com/cloudogu/k8s-component-lib/api/v1"
	"sigs.k8s.io/controller-runtime/pkg/config"

	"github.com/go-logr/logr"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

type mockDefinition struct {
	Arguments   []interface{}
	ReturnValue interface{}
}

func getCopyMap(definitions map[string]mockDefinition) map[string]mockDefinition {
	newCopyMap := map[string]mockDefinition{}
	for k, v := range definitions {
		newCopyMap[k] = v
	}
	return newCopyMap
}

func getNewMockManager(t *testing.T, expectedErrorOnNewManager error, definitions map[string]mockDefinition) manager.Manager {
	k8sManager := NewMockManager(t)
	ctrl.NewManager = func(config *rest.Config, options manager.Options) (manager.Manager, error) {
		for key, value := range definitions {
			k8sManager.Mock.On(key, value.Arguments...).Return(value.ReturnValue)
		}
		return k8sManager, expectedErrorOnNewManager
	}
	ctrl.SetLogger = func(l logr.Logger) {
		k8sManager.Mock.On("GetLogger").Return(l).Maybe()
	}

	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

	return k8sManager
}

func Test_startDoguOperator(t *testing.T) {
	// override default controller method to create a new manager
	oldNewManagerDelegate := ctrl.NewManager
	defer func() { ctrl.NewManager = oldNewManagerDelegate }()

	// override default controller method to retrieve a kube config
	oldGetConfigOrDieDelegate := ctrl.GetConfigOrDie
	defer func() { ctrl.GetConfigOrDie = oldGetConfigOrDieDelegate }()
	ctrl.GetConfigOrDie = func() *rest.Config {
		return &rest.Config{}
	}

	// override default controller method to retrieve a kube config
	oldGetConfigDelegate := ctrl.GetConfig
	defer func() { ctrl.GetConfig = oldGetConfigDelegate }()
	ctrl.GetConfig = func() (*rest.Config, error) {
		return &rest.Config{}, nil
	}

	// override default controller method to signal the setup handler
	oldHandler := ctrl.SetupSignalHandler
	defer func() { ctrl.SetupSignalHandler = oldHandler }()
	ctrl.SetupSignalHandler = func() context.Context {
		return context.TODO()
	}

	// override default controller method to retrieve a kube config
	oldSetLoggerDelegate := ctrl.SetLogger
	defer func() { ctrl.SetLogger = oldSetLoggerDelegate }()

	scheme := runtime.NewScheme()
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(v1.AddToScheme(scheme))
	myConfig := &rest.Config{}

	skipNameValidation := true
	defaultMockDefinitions := map[string]mockDefinition{
		"GetScheme":           {ReturnValue: scheme},
		"GetConfig":           {ReturnValue: myConfig},
		"GetEventRecorderFor": {Arguments: []interface{}{mock.Anything}, ReturnValue: nil},
		"GetControllerOptions": {ReturnValue: config.Controller{
			SkipNameValidation: &skipNameValidation,
		}},
		"Add":             {Arguments: []interface{}{mock.Anything}, ReturnValue: nil},
		"GetCache":        {ReturnValue: nil},
		"AddHealthzCheck": {Arguments: []interface{}{mock.Anything, mock.Anything}, ReturnValue: nil},
		"AddReadyzCheck":  {Arguments: []interface{}{mock.Anything, mock.Anything}, ReturnValue: nil},
		"Start":           {Arguments: []interface{}{mock.Anything}, ReturnValue: nil},
	}

	t.Run("Error on missing namespace environment variable", func(t *testing.T) {
		// given
		_ = os.Unsetenv("NAMESPACE")
		getNewMockManager(t, nil, defaultMockDefinitions)

		// when
		err := startOperator()

		// then
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to create new operator configuration: failed to read namespace: environment variable NAMESPACE must be set")
	})

	t.Setenv("NAMESPACE", "mynamespace")
	t.Setenv("RUNTIME", "local")
	t.Run("Test without logger environment variables", func(t *testing.T) {
		// given
		k8sManager := getNewMockManager(t, nil, defaultMockDefinitions)

		// when
		err := startOperator()

		// then
		require.NoError(t, err)
		mock.AssertExpectationsForObjects(t, k8sManager)
	})

	t.Run("Test with error on manager creation", func(t *testing.T) {
		// given
		expectedError := fmt.Errorf("this is my expected error")
		getNewMockManager(t, expectedError, map[string]mockDefinition{})

		// when
		err := startOperator()

		// then
		require.ErrorIs(t, err, expectedError)
	})

	t.Run(fmt.Sprintf("fail setup when error on %s", "Add"), func(t *testing.T) {
		// given
		expectedError := fmt.Errorf("this is my expected error for Add")
		adaptedMockDefinitions := getCopyMap(defaultMockDefinitions)
		adaptedMockDefinitions["Add"] = mockDefinition{
			Arguments:   adaptedMockDefinitions["Add"].Arguments,
			ReturnValue: expectedError,
		}
		delete(adaptedMockDefinitions, "AddHealthzCheck")
		delete(adaptedMockDefinitions, "AddReadyzCheck")
		delete(adaptedMockDefinitions, "Start")

		getNewMockManager(t, nil, adaptedMockDefinitions)

		// when
		err := startOperator()

		// then
		require.ErrorIs(t, err, expectedError)
	})
	t.Run(fmt.Sprintf("fail setup when error on %s", "AddHealthzCheck"), func(t *testing.T) {
		// given
		expectedError := fmt.Errorf("this is my expected error for AddHealthzCheck")
		adaptedMockDefinitions := getCopyMap(defaultMockDefinitions)
		adaptedMockDefinitions["AddHealthzCheck"] = mockDefinition{
			Arguments:   adaptedMockDefinitions["AddHealthzCheck"].Arguments,
			ReturnValue: expectedError,
		}
		delete(adaptedMockDefinitions, "AddReadyzCheck")
		delete(adaptedMockDefinitions, "Start")
		getNewMockManager(t, nil, adaptedMockDefinitions)

		// when
		err := startOperator()

		// then
		require.ErrorIs(t, err, expectedError)
	})
	t.Run(fmt.Sprintf("fail setup when error on %s", "AddReadyzCheck"), func(t *testing.T) {
		// given
		expectedError := fmt.Errorf("this is my expected error for AddReadyzCheck")
		adaptedMockDefinitions := getCopyMap(defaultMockDefinitions)
		adaptedMockDefinitions["AddReadyzCheck"] = mockDefinition{
			Arguments:   adaptedMockDefinitions["AddReadyzCheck"].Arguments,
			ReturnValue: expectedError,
		}
		delete(adaptedMockDefinitions, "Start")
		getNewMockManager(t, nil, adaptedMockDefinitions)

		// when
		err := startOperator()

		// then
		require.ErrorIs(t, err, expectedError)
	})
	t.Run(fmt.Sprintf("fail setup when error on %s", "Start"), func(t *testing.T) {
		// given
		expectedError := fmt.Errorf("this is my expected error for Start")
		adaptedMockDefinitions := getCopyMap(defaultMockDefinitions)
		adaptedMockDefinitions["Start"] = mockDefinition{
			Arguments:   adaptedMockDefinitions["Start"].Arguments,
			ReturnValue: expectedError,
		}
		getNewMockManager(t, nil, adaptedMockDefinitions)

		// when
		err := startOperator()

		// then
		require.ErrorIs(t, err, expectedError)
	})
}
