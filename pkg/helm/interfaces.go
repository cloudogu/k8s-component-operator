package helm

import (
	"context"

	k8sv1 "github.com/cloudogu/k8s-component-lib/api/v1"
)

type configMapRefReader interface {
	GetValues(ctx context.Context, configMapReference *k8sv1.Reference) (string, error)
}
