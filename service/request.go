package service

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/fatih/structs"
	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v2"
)

type RequestBody interface {
	ContentType() string
	Interface() interface{}
	Raw() ([]byte, error)
	Bind(interface{}) error
}

type Request interface {
	Context() context.Context
	Service() string
	Handler() string
	Method() string
	Body() RequestBody
	Headers() map[string]string
}

type RequestOptions struct {
	Context context.Context
	Service string
	Handler string
	Method  string
	Body    RequestBody
	Headers map[string]string
}

type RequestOption func(*RequestOptions)
type RequestExtension func(Request) Request

func WithContext(ctx context.Context) RequestOption {
	return func(ro *RequestOptions) {
		ro.Context = ctx
	}
}

func WithBody(ctype string, body interface{}) RequestOption {
	return func(ro *RequestOptions) {
		ro.Body = NewRequestBody(ctype, body)
	}
}

func WithHeaders(hdrs map[string]string) RequestOption {
	return func(ro *RequestOptions) {
		ro.Headers = hdrs
	}
}

type request struct {
	ctx     context.Context
	headers map[string]string
	body    RequestBody

	service string
	handler string
	method  string
}

func NewRequest(ctx context.Context, svc, hndlr, meth string, hdrs map[string]string, body RequestBody) Request {
	return &request{
		ctx:     ctx,
		headers: hdrs,
		body:    body,
		service: svc,
		handler: hndlr,
		method:  meth,
	}
}

func newOptRequest(opts RequestOptions) Request {
	return NewRequest(
		opts.Context,
		opts.Service,
		opts.Handler,
		opts.Method,
		opts.Headers,
		opts.Body,
	)
}

func NewRequestFromOptions(opts ...RequestOption) Request {
	reqOpts := RequestOptions{}
	for _, o := range opts {
		o(&reqOpts)
	}

	return newOptRequest(reqOpts)
}

func (r *request) Context() context.Context {
	return r.ctx
}

func (r *request) Service() string {
	return r.service
}

func (r *request) Handler() string {
	return r.handler
}

func (r *request) Method() string {
	return r.method
}

func (r *request) Body() RequestBody {
	return r.body
}

func (r *request) Headers() map[string]string {
	return r.headers
}

type reqBody struct {
	contentType string
	body        interface{}
}

func NewRequestBody(ct string, body interface{}) RequestBody {
	return &reqBody{
		contentType: ct,
		body:        body,
	}
}

func (rb *reqBody) ContentType() string {
	return rb.contentType
}

func (rb *reqBody) Interface() interface{} {
	return rb.body
}

func (rb *reqBody) Raw() ([]byte, error) {
	raw, ok := rb.body.([]byte)
	if !ok {
		return nil, errors.New("request body cannot be returned as []byte")
	}
	return raw, nil
}

func (rb *reqBody) Bind(v interface{}) error {
	switch rb.contentType {
	case "application/json":
		raw, err := rb.Raw()
		if err != nil {
			return err
		}
		return json.Unmarshal(raw, v)
	case "application/yaml":
		raw, err := rb.Raw()
		if err != nil {
			return err
		}
		return yaml.Unmarshal(raw, v)
	default:
		bm := structs.Map(rb.body)
		if bm == nil {
			return errors.New("could not bind request body")
		}
		return mapstructure.Decode(bm, v)
	}
}
