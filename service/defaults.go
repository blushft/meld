package service

import (
	"github.com/blushft/meld/service/handler"
	"github.com/blushft/meld/service/handler/rpc"
	"github.com/blushft/meld/service/method"
)

var (
	DefaultMethods  = map[string]method.MethodInvocation{"rpc": method.NewMethod}
	DefaultHandlers = map[string]handler.HandlerInvocation{"rpc": rpc.NewHandler}
)
