package yaml

import (
	"fmt"
	"io"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

type documentSplitter struct {
	decoder         decoder
	err             error
	currentDocument *runtime.RawExtension
}

func NewDocumentSplitter() DocumentSplitter {
	return &documentSplitter{}
}

func (s *documentSplitter) WithReader(r io.Reader) DocumentSplitter {
	s.decoder = yaml.NewYAMLOrJSONDecoder(r, 100)
	return s
}

func (s *documentSplitter) Next() bool {
	if err := s.decoder.Decode(s.currentDocument); err != nil {
		s.err = fmt.Errorf("failed to decode next yaml document: %w", err)
		s.currentDocument = nil
		return false
	}

	return true
}

func (s *documentSplitter) Err() error {
	if s.err == io.EOF {
		return nil
	}
	return s.err
}

func (s *documentSplitter) Bytes() []byte {
	if s.currentDocument != nil {
		return s.currentDocument.Raw
	}

	return nil
}

func (s *documentSplitter) Object() runtime.Object {
	if s.currentDocument != nil {
		return s.currentDocument.Object
	}

	return nil
}
