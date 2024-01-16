package yaml

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"testing"
)

func Test_documentSplitter_WithReader(t *testing.T) {
	sut := &documentSplitter{}

	withReader := sut.WithReader(new(bytes.Buffer))
	assert.NotEmpty(t, sut.decoder)
	assert.Same(t, sut, withReader)
}

func TestNewDocumentSplitter(t *testing.T) {
	newSplitter := NewDocumentSplitter()
	assert.NotNil(t, newSplitter)
}

func Test_documentSplitter_Err(t *testing.T) {
	tests := []struct {
		name    string
		err     error
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "should return nil if error is nil",
			err:     nil,
			wantErr: assert.NoError,
		},
		{
			name:    "should return nil if error is EOF",
			err:     io.EOF,
			wantErr: assert.NoError,
		},
		{
			name: "should return any other error",
			err:  assert.AnError,
			wantErr: func(t assert.TestingT, err error, _ ...interface{}) bool {
				return assert.ErrorIs(t, err, assert.AnError)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &documentSplitter{err: tt.err}
			tt.wantErr(t, s.Err())
		})
	}
}

func Test_documentSplitter_Bytes(t *testing.T) {
	tests := []struct {
		name            string
		currentDocument *runtime.RawExtension
		want            []byte
	}{
		{
			name:            "should return nil if document is nil",
			currentDocument: nil,
			want:            nil,
		},
		{
			name:            "should return document bytes if document is not nil",
			currentDocument: &runtime.RawExtension{Raw: []byte("banana")},
			want:            []byte("banana"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &documentSplitter{currentDocument: tt.currentDocument}
			assert.Equal(t, tt.want, s.Bytes())
		})
	}
}

func Test_documentSplitter_Object(t *testing.T) {
	tests := []struct {
		name            string
		currentDocument *runtime.RawExtension
		want            runtime.Object
	}{
		{
			name:            "should return nil if document is nil",
			currentDocument: nil,
			want:            nil,
		},
		{
			name:            "should return object if document is not nil",
			currentDocument: &runtime.RawExtension{Object: &v1.ConfigMap{}},
			want:            &v1.ConfigMap{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &documentSplitter{currentDocument: tt.currentDocument}
			assert.Equal(t, tt.want, s.Object())
		})
	}
}

func Test_documentSplitter_Next(t *testing.T) {
	t.Run("should return false on EOF", func(t *testing.T) {
		// given
		decoderMock := newMockDecoder(t)
		sut := &documentSplitter{decoder: decoderMock}
		decoderMock.EXPECT().Decode(&runtime.RawExtension{}).Return(io.EOF)

		// when
		actual := sut.Next()

		// then
		assert.ErrorIs(t, sut.err, io.EOF)
		assert.False(t, actual)
	})
	t.Run("should return false on other errors", func(t *testing.T) {
		// given
		decoderMock := newMockDecoder(t)
		sut := &documentSplitter{decoder: decoderMock}
		decoderMock.EXPECT().Decode(&runtime.RawExtension{}).Return(assert.AnError)

		// when
		actual := sut.Next()

		// then
		assert.ErrorIs(t, sut.err, assert.AnError)
		assert.False(t, actual)
	})
	t.Run("should return true if no errors", func(t *testing.T) {
		// given
		decoderMock := newMockDecoder(t)
		sut := &documentSplitter{decoder: decoderMock}
		decoderMock.EXPECT().Decode(&runtime.RawExtension{}).Return(nil)

		// when
		actual := sut.Next()

		// then
		assert.NoError(t, sut.err)
		assert.True(t, actual)
	})
}
