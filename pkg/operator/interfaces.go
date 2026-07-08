package operator

import (
	component "github.com/cloudogu/k8s-component-lib/client"
)

type componentClient interface {
	component.ComponentInterface
}

type helmClient interface {
	MarkReleaseAsFailed(name string, description string) error
}
