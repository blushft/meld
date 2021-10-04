package endpoint

import (
	"context"
	"fmt"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

type Tester struct {
	StringVal string
}

type TesterRequest struct {
	Value string
}

type TesterResponse struct {
	Result string
}

func (t *Tester) Test(ctx context.Context, req TesterRequest, resp *TesterResponse) error {
	spew.Dump("calling with args", ctx, req, resp)

	resp.Result = fmt.Sprintf("%s, %s", t.StringVal, req.Value)

	return nil
}

func TestCallEndpoint(t *testing.T) {
	tst := &Tester{StringVal: "test run"}
	eps, err := Extract(tst)
	if err != nil {
		t.Fatal(err)
	}

	spew.Dump(eps)

	req := NewRequest(context.Background(), TesterRequest{Value: "Test 1"})
	resp := &TesterResponse{}
	if err := eps[0].Call(req, resp); err != nil {
		t.Fatal(err)
	}

	spew.Dump(resp)
}
