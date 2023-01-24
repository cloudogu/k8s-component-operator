package helm

import (
	"fmt"
	helmclient "github.com/mittwald/go-helm-client"
	"helm.sh/helm/v3/pkg/action"
	ctrl "sigs.k8s.io/controller-runtime"
)

const (
	helmRepositoryCache  = "/tmp/.helmcache"
	helmRepositoryConfig = "/tmp/.helmrepo"
)

// New create a new instance of the helm client.
func New(namespace string, debug bool, debugLog action.DebugLog) (helmclient.Client, error) {
	opt := &helmclient.RestConfClientOptions{
		Options: &helmclient.Options{
			Namespace:        namespace,
			RepositoryCache:  helmRepositoryCache,
			RepositoryConfig: helmRepositoryConfig,
			Debug:            debug,
			DebugLog:         debugLog,
			Linting:          true,
		},
		RestConfig: ctrl.GetConfigOrDie(),
	}

	helmClient, err := helmclient.NewClientFromRestConf(opt)
	if err != nil {
		return nil, fmt.Errorf("failed to create helm client: %w", err)
	}

	return helmClient, nil
}
