package helm

import (
	"context"

	k8sv1 "github.com/cloudogu/k8s-component-lib/api/v1"
	"github.com/cloudogu/k8s-component-operator/pkg/yaml"
)

type configMapRefReader interface {
	GetValues(ctx context.Context, configMapReference *k8sv1.Reference) (string, error)
}

//nolint:unused
//goland:noinspection GoUnusedType
type yamlSerializer interface {
	yaml.Serializer
}
