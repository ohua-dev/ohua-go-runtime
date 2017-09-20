package backend

import (
	"reflect"
)

type CallableSfn struct{
	t interface{}
	sfn string
	id int
}

type LinkedDep struct{
	source int
	sourceIdx int
	target int
	targetIdx int
}

type RuntimeGraph struct {
	sfns []CallableSfn
	deps []LinkedDep
}

//type myType struct {
//	s int
//}
//
//func (m myType) myFunc(b int) {
//	m.s = b
//}

func (graph RuntimeGraph) Exec() interface{} {
	// TODO execute the graph here: build the operators, arcs and goroutines

	// create lambdas for each of the stateful functions and store them somewhere
	f := func(args ... interface{}) interface{} {
		// FIXME instead of casting back and forth we might leave the values just as they are and deal only with reflected values.
		var reflected_args []reflect.Value
		for _, arg := range args {
			reflected_args = append(reflected_args, reflect.ValueOf(arg))
		}
		// assumes only a single return result
		result := reflect.ValueOf(graph.sfns[0].t).MethodByName(graph.sfns[0].sfn).Call(reflected_args)
		return result[0].Interface()
	}
	return nil
}
