package greeter

import (
	"context"
	"fmt"

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
