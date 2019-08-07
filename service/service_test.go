package service

import (
	"context"
	"fmt"
	"testing"
)

type Greeter struct{}

type HelloReq struct {
	Name string `json:"name,omitempty"`
}

type HelloResp struct {
	Message string `json:"message,omitempty"`
}

func (g *Greeter) Hello(ctx context.Context, req HelloReq, resp *HelloResp, opts ...HandlerOption) error {
	resp.Message = fmt.Sprintf("Hello, %s", req.Name)
	return nil
}

type WelcomeReq struct {
	Salutory string `json:"salutory,omitempty"`
	Name     string `json:"name,omitempty"`
}

type WelcomeResp struct {
	Message string `json:"message,omitempty"`
}

func (g *Greeter) Welcome(ctx context.Context, req WelcomeReq, resp *WelcomeResp, opts ...HandlerOption) error {
	s := req.Name
	if req.Salutory != "" {
		s = fmt.Sprintf("%s %s", req.Salutory, req.Name)
	}
	resp.Message = fmt.Sprintf("Welcome to meld, %s", s)
	return nil
}

type StatusCheck struct{}

func (s *StatusCheck) Check(ctx context.Context, req struct{}, resp *string, opts ...HandlerOption) error {
	*resp = "ok"
	return nil
}

func Test_newService(t *testing.T) {
	type args struct {
		opts []Option
		h    interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test",
			args: args{
				opts: []Option{
					Name("testsvc"),
					Namespace("meld.test"),
					WithTag("test", "testing", "greet"),
					WithLabel("release", "latest"),
					Version("0.1.0"),
				},
				h: &Greeter{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newService(tt.args.h, tt.args.opts...)
			h := NewHandler(&StatusCheck{})
			s.Handle(h)
			fmt.Println(s.Usage())
			fmt.Println(s.Handlers())

		})
	}
}

func Test_service_Call(t *testing.T) {
	type args struct {
		req  Request
		resp *HelloResp
	}
	tests := []struct {
		name    string
		s       Service
		args    args
		wantErr bool
	}{
		{
			name: "test",
			s: newService(&Greeter{}, []Option{
				Name("testsvc"),
				Namespace("meld.test"),
				WithTag("test", "testing", "greet"),
				WithLabel("release", "latest"),
				Version("0.1.0"),
			}...),
			args: args{
				req: NewRequest(
					context.Background(),
					"testsvc",
					"Greeter",
					"Hello",
					nil,
					NewRequestBody("struct", HelloReq{Name: "Tester"}),
				),
				resp: &HelloResp{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.Call(tt.args.req, tt.args.resp)
			fmt.Println(tt.args.resp.Message)
		})
	}
}
