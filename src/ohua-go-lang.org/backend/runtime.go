package backend

import (
	"reflect"
	"sort"
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

type BySrcIndex []LinkedDep
// feels to me like the first two functions are always the same for arrays!!
func (b BySrcIndex) Len() int { return len(b) }
func (a BySrcIndex) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a BySrcIndex) Less(i, j int) bool { return a[i].sourceIdx < a[j].sourceIdx }


func (graph RuntimeGraph) Exec() interface{} {
	// ops first
	var ops map[int]Operator
	ops = make(map[int]Operator)
	for _, sfn := range graph.sfns {
		op := Operator{}
		op.f = func(args ... interface{}) interface{} {
			// FIXME instead of casting back and forth we might leave the values just as they are and deal only with reflected values.
			var reflected_args []reflect.Value
			for _, arg := range args {
				reflected_args = append(reflected_args, reflect.ValueOf(arg))
			}

			// FIXME need to handle destructuring properly! -> nth in the language!
			// assumes only a single return result
			result := reflect.ValueOf(graph.sfns[0].t).MethodByName(graph.sfns[0].sfn).Call(reflected_args)
			return result[0].Interface()
		}
		ops[sfn.id] = op
	}

	// sort deps
	deps := make(map[int][]LinkedDep)
	for _, dep := range graph.deps {
		_, exists := deps[dep.source]
		if  !exists {
			deps[dep.source] = make([]LinkedDep, 0)
		}
		deps[dep.source] = append(deps[dep.source], dep)
	}

	for _, value := range deps { sort.Sort(BySrcIndex(value)) }

	// TODO intersperse environment args

	// now create the arcs as channels
	for _, local_deps := range deps {
		for _, dep := range local_deps {
			c := make(chan interface{})
			source := ops[dep.source]
			source.out_channels = append(source.out_channels, c)
			target := ops[dep.target]
			target.in_channels = append(target.in_channels, c)
		}
	}

	// kick off the execution
	for _, op := range ops { go op.Get_exec() }

	// TODO wait for the last go routine to finish -> define a channel for the last op

	// TODO return the final result (if there is any)
	return nil
}
