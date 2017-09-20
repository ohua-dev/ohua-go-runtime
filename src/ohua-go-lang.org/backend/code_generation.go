package backend

import (
	"encoding/json"
	"os"
	"fmt"
	"strings"
	"ohua-go-lang.org/frontend"
	"reflect"
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
	main_template = "import (_ \"ohua-go-lang.org/backend\"\n " +
		"%s " + // import code
			")\n\n\n" +
			"func exec_ohua() interface{} {\n" +
				"var graph backend.RuntimeGraph\n" +
		 		"graph = RuntimeGraph{sfns: make([]CallableSfn, 0), deps: make([]LinkedDep, 0)}\n" +
				"%s\n" +
		        "\n" +
		        "%s\n\n" +
				"return graph.Exec()}"

	import_template := func(pkg_name string) string {
		return fmt.Sprintf("\t \"%s\"", pkg_name)
	}

	creation_template := func (var_id int, sfn_type string, idx int, sfn_id int, fn_name string) string {
		return fmt.Sprintf("pVar_%i := %s{}\n" +
							"graph.sfns[%i] = append(graph.sfns[%i], CallableSfn{f: pVar_%i, sfn: %s, id: %i})\n",
								var_id, sfn_type, idx, var_id, sfn_id, fn_name)
	}

	deps_template := func (idx int, source int, src_idx int, target int, target_idx int) string {
		return fmt.Sprintf("graph.deps[%i] = append(graph.deps[%i], LinkedDep{source: %i, sourceIdx: %i, target: %i, targetIdx: %i})",
			idx, source, src_idx, target, target_idx)
	}

	var import_code []string
	var sfn_code []string
	for idx, sfn := range graph.sfns  {
		x := strings.Split(sfn.t, ".")
		pkg_name := x[0]
		fn_name := x[1]

		import_code = append(import_code, import_template(pkg_name))
		sfn_code = append(sfn_code, creation_template(idx, sfn.t, idx, sfn.id, fn_name))
	}

	var deps_code []string
	for idx, dep := range graph.deps {
		deps_code = append(deps_code, deps_template(idx, dep.source.sfn, dep.source.index, dep.target.sfn, dep.target.index))
	}

	var import_code_complete string
	var sfn_code_complete string
	var deps_code_complete string
	import_code_complete = strings.Join(import_code, "")
	sfn_code_complete = strings.Join(sfn_code, "")
	deps_code_complete = strings.Join(deps_code, "")

	var final_code string
	final_code = fmt.Sprintf(main_template, import_code_complete, sfn_code_complete, deps_code_complete)

	f, err := os.Create("gen_exec.go")// TODO decide on a proper directory
	check(err)
	defer f.Close()

	_, err1 := f.WriteString(final_code)
	check(err1)

	// TODO it is yet unclear how can compile this newly generated file without taking the detour trough the os.
}