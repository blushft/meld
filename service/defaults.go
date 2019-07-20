package service

import (
	"github.com/blushft/meld/service/handler"
	"github.com/blushft/meld/service/handler/method"
)

var (
	DefaultMethods  = map[string]method.MethodInvocation{"rpc": method.NewMethod}
	DefaultHandlers = map[string]handler.HandlerInvocation{"rpc": handler.NewHandler}
)
