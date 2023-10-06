package v1

import (
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"reflect"
	"testing"
)

func TestComponent_GetHelmChartSpec(t *testing.T) {
	type fields struct {
		TypeMeta   v1.TypeMeta
		ObjectMeta v1.ObjectMeta
		Spec       ComponentSpec
		Status     ComponentStatus
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{name: "should use deployNamespace if specified", fields: fields{Spec: ComponentSpec{DeployNamespace: "longhorn"}}, want: "longhorn"},
		{name: "should use regular namespace if no deployNamespace if specified", fields: fields{ObjectMeta: v1.ObjectMeta{Namespace: "ecosystem"}, Spec: ComponentSpec{DeployNamespace: ""}}, want: "ecosystem"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Component{
				TypeMeta:   tt.fields.TypeMeta,
				ObjectMeta: tt.fields.ObjectMeta,
				Spec:       tt.fields.Spec,
				Status:     tt.fields.Status,
			}
			if got := c.GetHelmChartSpec().Namespace; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetHelmChartSpec() = %v, want %v", got, tt.want)
			}
		})
	}
}
