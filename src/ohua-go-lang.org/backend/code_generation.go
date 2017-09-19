package backend

import (
	"encoding/json"
	"os"
	"fmt"
	"strings"
)

/*
  This file contains the final steps of our compilation process.
  It gets a JSON string must turn it into a GO program that uses the functionality in runtime.
 */

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type ref struct{
	sfn int
	index int
}

type arc struct{
	source ref
	target ref
}

type sfn struct{
	id int
	t string
}

type compileDfGraph struct{
	// TODO mimic the compilation graph for easy JSON-to-GO conversion
	sfns []sfn `json:operators`
	deps []arc `json:arcs`
}

func genCode(df_graph string){
	// TODO JSON-to-GO
	var graph compileDfGraph

	var main_template string
	main_template = "import \"ohua-go-lang.org/backend\" \n\n\n" +
			"func exec_ohua(){ \n" +
				"var graph backend.RuntimeGraph\n" +
		 		"graph = RuntimeGraph{}" + // TODO fill the graph structure here
				"%s\n" +
				"return graph.exec()}"
	var creation_template string
	// TODO refine
	creation_template = "pVar_%s := %s{}\n" +
						"graph[%i] = pVar_%s\n"

	var init_code []string

	for _, sfn := range graph.sfns  {
		init_code = append(init_code, fmt.Sprintf(creation_template, sfn.id, sfn.t))
	}

	// TODO capture arc info into the graph as well

	var init_code_complete string
	init_code_complete = strings.Join(init_code, "")

	var final_code string
	final_code = fmt.Sprintf(main_template, init_code_complete)

	f, err := os.Create("gen_exec.go")// TODO decide on a proper directory
	check(err)
	defer f.Close()

	_, err1 := f.WriteString(final_code)
	check(err1)

	// TODO it is yet unclear how can compile this newly generated file without taking the detour trough the os.
}