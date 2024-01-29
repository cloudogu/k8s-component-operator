package yaml

import (
	_ "embed"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

//go:embed testdata/map.yaml
var mapYaml string

//go:embed testdata/invalid.yaml
var invalidYamlBytes []byte

func TestNewSerializer(t *testing.T) {
	newSerializer := NewSerializer()
	assert.NotNil(t, newSerializer)
}

func Test_serializer_Marshal(t *testing.T) {
	tests := []struct {
		name    string
		o       interface{}
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "should fail to marshal map with interface keys",
			o:       map[interface{}]string{1: "banana", 2: "apple"},
			want:    "",
			wantErr: assert.Error,
		},
		{
			name:    "should succeed to marshal map",
			o:       map[string]string{"a": "b", "c": "d"},
			want:    mapYaml,
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &serializer{}
			got, err := s.Marshal(tt.o)
			if !tt.wantErr(t, err, fmt.Sprintf("Marshal(%v)", tt.o)) {
				return
			}
			assert.Equalf(t, tt.want, string(got), "Marshal(%v)", tt.o)
		})
	}
}

func Test_serializer_Unmarshal(t *testing.T) {
	type args struct {
		y []byte
		o interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantObj interface{}
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "should fail for invalid yaml",
			args: args{
				y: invalidYamlBytes,
				o: &map[string]interface{}{},
			},
			wantObj: &map[string]interface{}{},
			wantErr: assert.Error,
		},
		{
			name: "should succeed for valid yaml",
			args: args{
				y: []byte(mapYaml),
				o: &map[string]interface{}{},
			},
			wantObj: &map[string]interface{}{"a": "b", "c": "d"},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &serializer{}
			tt.wantErr(t, s.Unmarshal(tt.args.y, tt.args.o), fmt.Sprintf("Unmarshal(%v, %v)", tt.args.y, tt.args.o))
			assert.Equal(t, tt.wantObj, tt.args.o)
		})
	}
}
