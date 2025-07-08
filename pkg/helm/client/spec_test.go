package client

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_haveSameKeyWithDifferentValues(t *testing.T) {
	tests := []struct {
		name string
		a    map[string]interface{}
		b    map[string]interface{}
		want bool
	}{
		{
			name: "two empty maps",
			a:    map[string]interface{}{},
			b:    map[string]interface{}{},
			want: false,
		},
		{
			name: "only different keys, more keys in a",
			a: map[string]interface{}{
				"a": "b",
				"c": "d",
			},
			b: map[string]interface{}{
				"e": "f",
			},
			want: false,
		},
		{
			a: map[string]interface{}{
				"e": "f",
			},
			name: "only different keys, more keys in b",
			b: map[string]interface{}{
				"a": "b",
				"c": "d",
			},
			want: false,
		},
		{
			name: "two equal maps",
			a: map[string]interface{}{
				"a": "b",
			},
			b: map[string]interface{}{
				"a": "b",
			},
			want: true,
		},
		{
			name: "two nested maps with different keys",
			a: map[string]interface{}{
				"a": map[string]interface{}{
					"c": "d",
				},
			},
			b: map[string]interface{}{
				"a": map[string]interface{}{
					"e": "f",
				},
			},
			want: false,
		},
		{
			name: "two nested maps with same keys and values",
			a: map[string]interface{}{
				"a": map[string]interface{}{
					"c": "d",
				},
			},
			b: map[string]interface{}{
				"a": map[string]interface{}{
					"c": "d",
				},
			},
			want: true,
		},
		{
			name: "two nested maps with same keys different values",
			a: map[string]interface{}{
				"a": map[string]interface{}{
					"c": "d",
				},
			},
			b: map[string]interface{}{
				"a": map[string]interface{}{
					"c": "e",
				},
			},
			want: true,
		},
		{
			name: "two nested maps with same keys different values",
			a: map[string]interface{}{
				"a": map[string]interface{}{
					"b": []map[string]interface{}{
						{
							"c": "d",
						},
					},
				},
			},
			b: map[string]interface{}{
				"a": map[string]interface{}{
					"b": []map[string]interface{}{
						{
							"e": "f",
						},
					},
				},
			},
			want: false,
		},
		{
			name: "second map overwrites first map array",
			a: map[string]interface{}{
				"a": map[string]interface{}{
					"b": []map[string]interface{}{
						{
							"c": "d",
						},
					},
				},
			},
			b: map[string]interface{}{
				"a": map[string]interface{}{
					"b": "e",
				},
			},
			want: true,
		},
		{
			name: "first map overwrites second map array",
			a: map[string]interface{}{
				"a": map[string]interface{}{
					"b": "e",
				},
			},
			b: map[string]interface{}{
				"a": map[string]interface{}{
					"b": []map[string]interface{}{
						{
							"c": "d",
						},
					},
				},
			},
			want: true,
		},
		{
			name: "both slices have multiple entries with one conflict",
			a: map[string]interface{}{
				"a": map[string]interface{}{
					"b": []map[string]interface{}{
						{
							"c": "d",
						},
						{
							"x": "y",
						},
						{
							"z": "a",
						},
					},
				},
			},
			b: map[string]interface{}{
				"a": map[string]interface{}{
					"b": []map[string]interface{}{
						{
							"y": "x",
						},
						{
							"a": "z",
						},
						{
							"c": "d",
						},
					},
				},
			},
			want: true,
		},
		{
			name: "both slices have multiple entries without conflicts",
			a: map[string]interface{}{
				"a": map[string]interface{}{
					"b": []map[string]interface{}{
						{
							"c": "d",
						},
						{
							"x": "y",
						},
						{
							"z": "a",
						},
					},
				},
			},
			b: map[string]interface{}{
				"a": map[string]interface{}{
					"b": []map[string]interface{}{
						{
							"y": "x",
						},
						{
							"a": "z",
						},
						{
							"g": "h",
						},
					},
				},
			},
			want: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equalf(t, test.want, hasSameValuesConfigured(test.a, test.b), "hasSameValuesConfigured(%v, %v)", test.a, test.b)
		})
	}
}
