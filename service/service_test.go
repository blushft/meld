package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/blushft/meld/service/method"
)

type Greeter struct{}

type HelloReq struct {
	Name string `json:"name,omitempty"`
}

type HelloResp struct {
	Message string `json:"message,omitempty"`
}

func (g *Greeter) Hello(ctx context.Context, req HelloReq, resp HelloResp, opts ...method.MethodOption) error {
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

func (g *Greeter) Welcome(ctx context.Context, req WelcomeReq, resp WelcomeResp, opts ...method.MethodOption) error {
	s := req.Name
	if req.Salutory != "" {
		s = fmt.Sprintf("%s %s", req.Salutory, req.Name)
	}
	resp.Message = fmt.Sprintf("Welcome to meld, %s", s)
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
			fmt.Println(s.Usage())
		})
	}
}
