package service

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/fatih/structs"
	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v2"
)

// RequestBody interface allows interaction with the raw request or Bind will
// map it onto a struct
type RequestBody interface {
	ContentType() string
	Interface() interface{}
	Raw() ([]byte, error)
	Bind(interface{}) error
}

// Request describes the Service Call and contains the Body and Headers
type Request interface {
	Context() context.Context
	Service() string
	Handler() string
	Method() string
	Body() RequestBody
	Headers() map[string]string
}

// RequestOptions allows for declaritive construction of a Request
type RequestOptions struct {
	Context context.Context
	Service string
	Handler string
	Method  string
	Body    RequestBody
	Headers map[string]string
}

// RequestOption acts as a setter for RequestOptions
type RequestOption func(*RequestOptions)

// RequestExtension is middleware for Requests
type RequestExtension func(Request) Request

// WithContext sets the Context of the Request
func WithContext(ctx context.Context) RequestOption {
	return func(ro *RequestOptions) {
		ro.Context = ctx
	}
}

// WithBody sets the body of the request
func WithBody(ctype string, body interface{}) RequestOption {
	return func(ro *RequestOptions) {
		ro.Body = NewRequestBody(ctype, body)
	}
}

// WithHeaders sets the Headers of the Request
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

// NewRequest returns a Request with functional argument construction
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

// NewRequestFromOptions returns a new Request by declaritive construction
func NewRequestFromOptions(opts ...RequestOption) Request {
	reqOpts := RequestOptions{}
	for _, o := range opts {
		o(&reqOpts)
	}

	return newOptRequest(reqOpts)
}

// Context returns the Context of the request
func (r *request) Context() context.Context {
	return r.ctx
}

// Service returns the name of the targeted Service in the Request
func (r *request) Service() string {
	return r.service
}

// Handler returns the name of the targeted Handler in the Request
func (r *request) Handler() string {
	return r.handler
}

// Method returns the name of the targeted Method in the Request
func (r *request) Method() string {
	return r.method
}

// Body returns the RequestBody
func (r *request) Body() RequestBody {
	return r.body
}

// Headers returns any Headers attached to the Request
func (r *request) Headers() map[string]string {
	return r.headers
}

type reqBody struct {
	contentType string
	body        interface{}
}

// NewRequestBody returns a RequestBody with a specified Content-Type and arbitrary Content
func NewRequestBody(ct string, body interface{}) RequestBody {
	return &reqBody{
		contentType: ct,
		body:        body,
	}
}

// ContentType returns the content type of the Body interface{}
func (rb *reqBody) ContentType() string {
	return rb.contentType
}

// Interface returns the Body as an interface{}
func (rb *reqBody) Interface() interface{} {
	return rb.body
}

// Raw returns the body as []byte or error if the raw data is not in that format
func (rb *reqBody) Raw() ([]byte, error) {
	raw, ok := rb.body.([]byte)
	if !ok {
		return nil, errors.New("request body cannot be returned as []byte")
	}
	return raw, nil
}

// Bind will attempt to Unmarshal an encoded Body to an interface{} using the defined
// Content-Type
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
