package handler

import (
	"context"
	"fmt"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

type Tester struct {
	text string
}

type TesterResp struct {
	Val string
}

func (t *Tester) TestMe(ctx context.Context, req string, resp *TesterResp, opts ...HandlerOption) error {
	spew.Dump("calling with args: ", ctx, req, resp)

	ctxVal := ctx.Value(000)
	ctxS, ok := ctxVal.(string)
	if !ok {
		ctxS = "no ctx val"
	}

	resp.Val = fmt.Sprintf("%s, %s = %s", req, t.text, ctxS)
	return nil
}

var (
	stringFn = func(ctx context.Context, req string, resp *string, opts ...string) error {
		*resp = "this is a string"
		return nil
	}
)

func Test_newRPCHandler(t *testing.T) {
	type args struct {
		v    interface{}
		opts []HandlerOption
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test_struct",
			args: args{
				v: &Tester{text: ", tester"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.WithValue(context.Background(), 000, "my ctx value")
			h := newRPCHandler(tt.args.v, tt.args.opts...)
			ret := &TesterResp{}
			tt.args.v.(*Tester).TestMe(ctx, "value", ret)
			fmt.Println("local return", ret.Val)
			cRet := &TesterResp{}
			h.Call(ctx, "TestMe", "valueX", cRet)
			fmt.Println(ret.Val)
		})
	}
}

func Test_newRPCHandlerFunc(t *testing.T) {
	type args struct {
		v    interface{}
		opts []HandlerOption
	}
	tests := []struct {
		name string
		args args
	}{

		{
			name: "test_func",
			args: args{
				v: stringFn,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//ctx := context.WithValue(context.Background(), 000, "my ctx value")
			h := newRPCHandler(tt.args.v, tt.args.opts...)
			spew.Dump(h)
		})
	}
}
