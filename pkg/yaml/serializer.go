package yaml

import (
	"sigs.k8s.io/yaml"
)

type serializer struct{}

func NewSerializer() Serializer {
	return &serializer{}
}

func (s *serializer) Marshal(o interface{}) ([]byte, error) {
	return yaml.Marshal(o)
}

func (s *serializer) Unmarshal(y []byte, o interface{}, opts ...yaml.JSONOpt) error {
	return yaml.Unmarshal(y, o, opts...)
}
