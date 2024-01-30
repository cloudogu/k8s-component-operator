package labels

import (
	"github.com/cloudogu/k8s-component-operator/pkg/yaml"
	"k8s.io/apimachinery/pkg/runtime"
)

type documentSplitter interface {
	yaml.DocumentSplitter
}

type unstructuredSerializer interface {
	runtime.Serializer
}

type unstructuredConverter interface {
	runtime.UnstructuredConverter
}

type genericYamlSerializer interface {
	yaml.Serializer
}
