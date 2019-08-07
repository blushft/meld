package service

import (
	"context"
	"reflect"
	"testing"
)

type (
	stringStr struct{}
)

func (str *stringStr) Stringer(ctx context.Context, req string, resp *string) error {
	*resp = "this is a string"
	return nil
}

var (
	mstringFn = func(ctx context.Context, req string, resp *string) error {
		*resp = "this is a string"
		return nil
	}
)

func Test_extractSig(t *testing.T) {
	type args struct {
		v reflect.Type
		d int
	}
	tests := []struct {
		name string
		args args
		want *MethodDef
	}{
		{
			name: "test_struct",
			args: args{
				v: reflect.TypeOf(&stringStr{}),
				d: 0,
			},
			want: &MethodDef{},
		},
		{
			name: "test_func",
			args: args{
				v: reflect.TypeOf(stringFn),
				d: 0,
			},
			want: &MethodDef{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := extractMethod(tt.args.v, tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("extractMethod() = %v, want %v", got, tt.want)
			}
		})
	}
}
