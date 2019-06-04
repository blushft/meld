package service

import "context"

type Request interface {
	Service() string
	Method() string
	Body() interface{}
	Options() *CallOptions
}

type CallOptions struct {
	Meta    map[string]string
	Context context.Context
}
