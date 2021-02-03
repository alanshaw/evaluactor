package evaluactor

import (
	"github.com/Shopify/go-lua"
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/cbor"
	"github.com/filecoin-project/go-state-types/exitcode"
	"github.com/filecoin-project/go-state-types/rt"
	"github.com/filecoin-project/lotus/chain/actors/builtin"
	"github.com/filecoin-project/specs-actors/v2/actors/runtime"
	"github.com/ipfs/go-cid"
	"github.com/multiformats/go-multihash"
)

//go:generate go run ./gen

// Actor is a Filecoin actor that can interpret lua sent to it in method params.
type Actor struct{}

// State is the stored state for this actor.
type State struct{}

const (
	_ = 0 // skip zero iota value; first usage of iota gets 1.
	// MethodEval is the numeric identifier of the method that evaluates lua and
	// returns the result.
	MethodEval = builtin.MethodConstructor + iota
)

// EvaluactorActorCodeCID is the CID by which this kind of actor will be identified.
var EvaluactorActorCodeCID = func() cid.Cid {
	builder := cid.V1Builder{Codec: cid.Raw, MhType: multihash.IDENTITY}
	c, err := builder.Sum([]byte("fil/1/eval"))
	if err != nil {
		panic(err)
	}
	return c
}()

// Address is the singleton address of this actor. Its value is 97
// (builtin.FirstNonSingletonActorId - 3), as 99 is reserved for the burnt funds
// singleton and 98 is the chaoas actor.
var Address = func() address.Address {
	// the address before the chaos actor (98)
	addr, err := address.NewIDAddress(97)
	if err != nil {
		panic(err)
	}
	return addr
}()

// ErrScriptRunFailure is the exit code associated with a runtime failure of the
// lua script.
const ErrScriptRunFailure = exitcode.FirstActorSpecificExitCode

// Exports defines the methods this actor exposes publicly.
func (a Actor) Exports() []interface{} {
	return []interface{}{
		builtin.MethodConstructor: a.Constructor,
		MethodEval:                a.Eval,
	}
}

// Code returns the CID of the actor.
func (a Actor) Code() cid.Cid { return EvaluactorActorCodeCID }

// State returns the nil state of the actor.
func (a Actor) State() cbor.Er { return new(State) }

// IsSingleton determines whether there are multiple instances of this actor.
func (a Actor) IsSingleton() bool { return true }

var _ rt.VMActor = Actor{}

// EvalParams are the parameters for the Eval method.
type EvalParams struct {
	Script string
}

// EvalReturn is the return value for the Eval method.
type EvalReturn struct {
	Result string
}

// Eval requests for this actor to evaluate some lua script and return the result.
func (a Actor) Eval(rt runtime.Runtime, params *EvalParams) *EvalReturn {
	rt.ValidateImmediateCallerAcceptAny()
	var res string
	l := lua.NewState()
	lua.OpenLibraries(l)

	_ = lua.NewMetaTable(l, "filMetaTable")
	lua.SetFunctions(l, []lua.RegistryFunction{{
		Name: "setresult",
		Function: func(l *lua.State) int {
			res = lua.CheckString(l, 1)
			return 0
		},
	}}, 0)
	l.SetGlobal("fil")
	lua.SetMetaTableNamed(l, "filMetaTable")

	err := lua.DoString(l, params.Script)
	if err != nil {
		rt.Abortf(ErrScriptRunFailure, "run script failure: %v", err)
	}
	return &EvalReturn{Result: res}
}

// Constructor will panic because the Evaluator actor is a singleton.
func (a Actor) Constructor(_ runtime.Runtime, _ *abi.EmptyValue) *abi.EmptyValue {
	panic("constructor should not be called; the Evaulactor actor is a singleton actor")
}
