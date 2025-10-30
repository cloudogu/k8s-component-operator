package configref

import v1 "k8s.io/client-go/kubernetes/typed/core/v1"

type configMapClient interface {
	v1.ConfigMapInterface
}
