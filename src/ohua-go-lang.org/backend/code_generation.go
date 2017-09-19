package backend

import (
	"encoding/json"
	"os"
	"fmt"
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

func gen_code(df_graph string){
	var main_template string
	main_template = "import \"ohua-go-lang.org/backend\" \n\n\n" +
			"func exec_ohua(){ \n" +
				"var graph backend.Df_graph\n" +
		 		"graph = Df_graph{}" + // TODO fill the graph structure here
				"%s\n" +
				"return graph.exec()}"
	var creation_template string
	creation_template = "%s "

	var instantiation_code string

	// TODO

	var final_code string
	final_code = fmt.Sprintf(main_template, instantiation_code)

	f, err := os.Create("gen_exec.go")// TODO decide on a proper directory
	check(err)
	defer f.Close()

	_, err1 := f.WriteString(final_code)
	check(err1)
}