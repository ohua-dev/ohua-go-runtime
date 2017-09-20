package backend

import (
	"reflect"
	"sort"
)

type sfn_call func(args ... interface{}) interface{}

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

type operator struct{
	f sfn_call
	in_channels []chan interface{}
	out_channels []chan interface{}
}

type BySrcIndex []LinkedDep
// feels to me like the first two functions are always the same for arrays!!
func (b BySrcIndex) Len() int { return len(b) }
func (a BySrcIndex) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a BySrcIndex) Less(i, j int) bool { return a[i].sourceIdx < a[j].sourceIdx }


func (graph RuntimeGraph) Exec() interface{} {
	// TODO execute the graph here: build the operators, arcs and goroutines

	// ops first
	var ops map[int]operator
	ops = make(map[int]operator)
	for _, sfn := range graph.sfns {
		op := operator{}
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

	// sort deps first
	deps := make(map[int][]LinkedDep)
	for _, dep := range graph.deps {
		_, exists := deps[dep.source]
		if  !exists {
			deps[dep.source] = make([]LinkedDep, 0)
		}
		deps[dep.source] = append(deps[dep.source], dep)
	}

	for _, value := range deps {
		sort.Sort(BySrcIndex(value))
	}

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

	// TODO kick off the goroutines

	// TODO return the final result (if there is any)
	return nil
}

func (op operator) exec_call_strict(){
	var in_vals []interface{}
	// strict call semantics
	for _, channel := range op.in_channels {
		in_vals = append(in_vals, <-channel)
	}

	// call the sfn
	ret_val := op.f(in_vals)

	// dispatch to all interested
	for _, channel := range op.out_channels {
		channel <- ret_val
	}
}
