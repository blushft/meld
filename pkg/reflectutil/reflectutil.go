package reflectutil

import (
	"reflect"
)

type Attributes struct {
	Type     string   `json:"type"`
	Value    string   `json:"value"`
	Kind     string   `json:"kind"`
	Pkg      string   `json:"pkg"`
	Name     string   `json:"name"`
	Variadic bool     `json:"variadic"`
	NumIn    int      `json:"num_in"`
	NumOut   int      `json:"num_out"`
	Args     []string `json:"args"`
	Returns  []string `json:"returns"`
}

func Attr(v interface{}) Attributes {
	attr := Attributes{}
	var vv reflect.Type

	if reflect.TypeOf(v).Kind() == reflect.Ptr {
		vv = reflect.TypeOf(v).Elem()
	} else {
		vv = reflect.TypeOf(v)
	}

	attr.Type = vv.String()
	attr.Value = vv.String()
	attr.Kind = vv.Kind().String()
	attr.Pkg = vv.PkgPath()
	attr.Name = vv.Name()

	if vv.Kind() == reflect.Func {
		attr.Variadic = vv.IsVariadic()
		attr.NumIn = vv.NumIn()
		attr.NumOut = vv.NumOut()

		for i := 0; i < vv.NumIn(); i++ {
			in := vv.In(i)
			kind := in.Kind().String()

			if in.Kind() == reflect.Ptr {
				kind = "*" + in.Elem().Kind().String()
			}

			if in.Kind() == reflect.Interface {
				kind = in.PkgPath() + "#" + in.Name()
			}

			attr.Args = append(attr.Args, kind)
		}
		for i := 0; i < vv.NumOut(); i++ {
			out := vv.Out(i)

			kind := out.Kind().String()
			if out.Kind() == reflect.Ptr {
				kind = "*" + out.Elem().Kind().String()
			}

			if out.Kind() == reflect.Interface {
				kind = out.PkgPath() + "#" + out.Name()
			}

			attr.Returns = append(attr.Returns, kind)
		}
	}

	return attr
}
