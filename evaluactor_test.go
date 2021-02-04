package evaluactor

import (
	"context"
	"testing"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	mock "github.com/filecoin-project/specs-actors/v2/support/mock"
	"github.com/ipfs/go-cid"
)

func TestEval(t *testing.T) {
	receiver := Address
	builder := mock.NewBuilder(context.Background(), receiver)

	rt := builder.Build(t)
	rt.ExpectValidateCallerAny()

	a := Actor{}
	src := `
	local fil = require("fil")
	fil.set_result("OMG")
`

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

func TestRuntimeCaller(t *testing.T) {
	caller, err := address.NewIDAddress(100)
	if err != nil {
		t.Fatal(err)
	}

	receiver := Address
	builder := mock.NewBuilder(context.Background(), receiver)
	builder.WithCaller(caller, cid.Undef)

	rt := builder.Build(t)
	rt.ExpectValidateCallerAny()

	a := Actor{}
	src := `
	local fil = require("fil")
	local caller = fil.runtime.caller()
	fil.set_result(caller)
`

	ret := rt.Call(a.Eval, &EvalParams{Script: src})
	eret, ok := ret.(*EvalReturn)
	if !ok {
		t.Fatal("invalid return value")
	}

	if eret.Result != caller.String() {
		t.Fatalf("expected result to be \"%s\", got: %v", caller.String(), eret.Result)
	}

	rt.Verify()
}

func TestRuntimeValueReceived(t *testing.T) {
	valueReceived := abi.NewTokenAmount(5000)

	receiver := Address
	builder := mock.NewBuilder(context.Background(), receiver)

	rt := builder.Build(t)
	rt.SetReceived(valueReceived)
	rt.ExpectValidateCallerAny()

	a := Actor{}
	src := `
	local fil = require("fil")
	local val = fil.runtime.value_received()
	fil.set_result(val)
`

	ret := rt.Call(a.Eval, &EvalParams{Script: src})
	eret, ok := ret.(*EvalReturn)
	if !ok {
		t.Fatal("invalid return value")
	}

	if eret.Result != valueReceived.String() {
		t.Fatalf("expected result to be \"%s\", got: %v", valueReceived.String(), eret.Result)
	}

	rt.Verify()
}
