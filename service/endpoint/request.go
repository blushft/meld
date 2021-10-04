package endpoint

import (
	"context"
)

type Request interface {
	Context() context.Context
	Body() interface{}
	Headers() Headers
}

type request struct {
	ctx     context.Context
	headers Headers
	body    interface{}
}

func NewRequest(ctx context.Context, v interface{}, headers ...map[string]string) Request {
	return &request{
		ctx:     ctx,
		headers: NewHeaders(headers...),
		body:    v,
	}
}

func (r *request) Context() context.Context {
	return r.ctx
}

func (r *request) Body() interface{} {
	return r.body
}

func (r *request) Headers() Headers {
	return r.headers
}
