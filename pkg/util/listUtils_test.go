package util

import (
	"fmt"
	"reflect"
	"strconv"
	"testing"
)

func TestMap(t *testing.T) {
	type args[T any, V any] struct {
		ts []T
		fn func(T) V
	}
	type testCase[T any, V any] struct {
		name string
		args args[T, V]
		want []V
	}
	tests := []testCase[int, string]{
		{
			name: "int to string",
			args: args[int, string]{
				ts: []int{1, 2, 3},
				fn: func(i int) string {
					return strconv.Itoa(i)
				},
			},
			want: []string{"1", "2", "3"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Map(tt.args.ts, tt.args.fn); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Map() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReduce(t *testing.T) {
	type args[T any, A any] struct {
		ts  []T
		acc A
		fn  func(value T, acc A) A
	}
	type testCase[T any, A any] struct {
		name string
		args args[T, A]
		want A
	}
	tests := []testCase[int, string]{
		{
			name: "should join all integers to a string seperated by a space character",
			args: args[int, string]{
				ts:  []int{0x1F61C, 0x1F92A, 0x1F92D},
				acc: "",
				fn: func(value int, acc string) string {
					utf8Emoji := string(rune(value))
					if acc == "" {
						return utf8Emoji
					}
					return fmt.Sprintf("%s %s", acc, utf8Emoji)
				},
			},
			want: "😜 🤪 🤭",
		},
		{
			name: "should list integers in binary format",
			args: args[int, string]{
				ts:  []int{0, 1, 2},
				acc: "Binary:",
				fn: func(value int, acc string) string {
					return fmt.Sprintf("%s\n- 0x%b", acc, value)
				},
			},
			want: "Binary:\n- 0x0\n- 0x1\n- 0x10",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Reduce(tt.args.ts, tt.args.acc, tt.args.fn); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Reduce() = %v, want %v", got, tt.want)
			}
		})
	}
}
