package function

import (
	"reflect"

	"github.com/blang/semver"
)

type Function interface {
	Name() string
	Request() *FuncDef
	Response() *FuncDef
	Options() *FuncOptions
}

type FuncDef struct {
	Name   string     `json:"name"`
	Type   string     `json:"type"`
	Values []*FuncDef `json:"values"`

	t reflect.Type
}

func (f *FuncDef) TypeOf() reflect.Type {
	return f.t
}

type FuncOptions struct {
	Name    string
	Version semver.Version
	Labels  map[string]string
	Tags    []string
}

type FuncOption func(*FuncOptions)

type FunctionInvocation func(rf reflect.Type, opts ...FuncOption) (Function, error)
