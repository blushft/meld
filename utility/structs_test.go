package utility

import (
	"reflect"
	"testing"
)

type TestStruct struct {
	Name string
	Port string
	Host string
	Path string
}

func TestStructs_Attr(t *testing.T) {
	type args struct {
		v interface{}
	}
	tests := []struct {
		name string
		s    *Structs
		args args
		want map[string]string
	}{
		{
			name: "test",
			s:    &Structs{},
			args: args{
				v: TestStruct{
					Name: "testobj",
					Port: "9990",
					Host: "",
					Path: "/test",
				},
			},
			want: map[string]string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Attr(tt.args.v); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Structs.Attr() = %v, want %+v", got, tt.want)
			}
		})
	}
}
