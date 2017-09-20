package backend

type sfn_call func(args ... interface{}) interface{}

type Operator struct{
	f sfn_call
	in_channels []chan interface{}
	out_channels []chan interface{}
}

func (op Operator) Exec_call_strict() bool {
	var in_vals []interface{}
	var done bool
	done = false
	// strict call semantics
	for _, channel := range op.in_channels {
		val, ok := <- channel
		if ok {
			if done {panic("Invariant broken! Data imbalance.")}
			in_vals = append(in_vals, val)
		} else {
			done = true
		}
	}

	// call the sfn
	ret_val := op.f(in_vals)

	// dispatch to all interested
	for _, channel := range op.out_channels {
		channel <- ret_val
	}

	return done
}

func (op Operator) shutdown(){
	for _, in := range op.out_channels {
		close(in)
	}
}

func (op Operator) Get_exec() (func ()){
	if len(op.in_channels) == 0 {
		return func() {
			op.Exec_call_strict()
			op.shutdown()
		}
	} else {
		return func() {
			for !op.Exec_call_strict() {}
			op.shutdown()
		}
	}
}
