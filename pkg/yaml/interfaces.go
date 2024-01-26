package yaml

import (
	"io"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/yaml"
)

type DocumentSplitter interface {
	WithReader(r io.Reader) DocumentSplitter
	Next() bool
	Err() error
	Bytes() []byte
	Object() runtime.Object
}

type Serializer interface {
	Marshal(o interface{}) ([]byte, error)
	Unmarshal(y []byte, o interface{}, opts ...yaml.JSONOpt) error
}

type decoder interface {
	Decode(into interface{}) error
}
