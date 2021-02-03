package main

import (
	"github.com/alanshaw/evaluactor"

	gen "github.com/whyrusleeping/cbor-gen"
)

func main() {
	if err := gen.WriteTupleEncodersToFile("./cbor_gen.go", "evaluactor",
		evaluactor.State{},
		evaluactor.EvalParams{},
		evaluactor.EvalReturn{},
	); err != nil {
		panic(err)
	}
}
