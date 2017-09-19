package backend

import "go/types"

type callableSfn struct{
	f func(args ... types.Object) // FIXME it must be something that we can call
	id int
}

type linkedDep struct{
	source *sfn
	sourceIdx int
	target *sfn
	targetIdx int
}

type RuntimeGraph struct {
	sfns []callableSfn
	deps []linkedDep
}

func (graph RuntimeGraph) Exec() types.Object {
	// TODO execute the graph here: build the operators, arcs and goroutines

	return nil
}
