package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/blushft/meld/service/handler"
)

type Greeter struct{}

type HelloReq struct {
	Name string `json:"name,omitempty"`
}

type HelloResp struct {
	Message string `json:"message,omitempty"`
}

func (g *Greeter) Hello(ctx context.Context, req HelloReq, resp *HelloResp, opts ...handler.HandlerOption) error {
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

func (g *Greeter) Welcome(ctx context.Context, req WelcomeReq, resp *WelcomeResp, opts ...handler.HandlerOption) error {
	s := req.Name
	if req.Salutory != "" {
		s = fmt.Sprintf("%s %s", req.Salutory, req.Name)
	}
	resp.Message = fmt.Sprintf("Welcome to meld, %s", s)
	return nil
}

type StatusCheck struct{}

func (s *StatusCheck) Check(ctx context.Context, req struct{}, resp *string, opts ...handler.HandlerOption) error {
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
			h, _ := handler.NewHandler(&StatusCheck{})
			s.Handle(h)
			fmt.Println(s.Usage())
			fmt.Println(s.Handlers())

		})
	}
}

func Test_service_Call(t *testing.T) {
	type args struct {
		ctx     context.Context
		handler string
		method  string
		req     interface{}
		resp    interface{}
		opts    []CallOptions
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
				ctx:     context.Background(),
				handler: "Greeter",
				method:  "Hello",
				req:     HelloReq{Name: "Tester"},
				resp:    &HelloResp{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.Call(tt.args.ctx, tt.args.handler, tt.args.method, tt.args.req, tt.args.resp, tt.args.opts...)
			fmt.Println(tt.args.resp.(*HelloResp).Message)
		})
	}
}
