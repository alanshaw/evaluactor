# evaluactor

A Filecoin actor that interprets lua.

## Install

```sh
go get github.com/alanshaw/evaluactor
```

## Usage

```go
package main

import (
  "fmt"
  "github.com/alanshaw/evaluactor"
  "github.com/filecoin-project/specs-actors/v2/actors/runtime"
)

func main () {
  var rt runtime.Runtime
  a := evaluactor.Actor{}
  src := `fil.setresult("OMG")` // TODO: replace with an interesting example 😝
  ret := a.Eval(rt, &evaluactor.EvalParams{Script:src})
  fmt.Println(ret.Result) // OMG
}
```
