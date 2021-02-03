package evaluactor

import (
	"context"
	"testing"

	mock "github.com/filecoin-project/specs-actors/v2/support/mock"
)

func TestEval(t *testing.T) {
	receiver := Address
	builder := mock.NewBuilder(context.Background(), receiver)

	rt := builder.Build(t)
	rt.ExpectValidateCallerAny()

	a := Actor{}
	src := `fil.setresult("OMG")`

	ret := rt.Call(a.Eval, &EvalParams{Script: src})
	eret, ok := ret.(*EvalReturn)
	if !ok {
		t.Fatal("invalid return value")
	}

	if eret.Result != "OMG" {
		t.Fatalf("expected result to be \"OMG\", got: %v", eret.Result)
	}

	rt.Verify()
}
