package frontend

import (
	"go/importer"
	"fmt"
	"go/types"
)

func Exists(pkg_name string, fn_name string) bool {
	pkg, err := importer.Default().Import(pkg_name)
	if err != nil {
		fmt.Printf("error: %s\n", err.Error())
		return false
	}else{
		for _, declName := range pkg.Scope().Names() {
			var t types.Object
			t = pkg.Scope().Lookup(declName)

			if t == nil { panic(fmt.Sprintf("Impossible!")) }
			fmt.Println(t)
			obj, _, _ := types.LookupFieldOrMethod(t.Type(), true, pkg, fn_name) // check for a function in a type
			if obj == nil {
				// continue
			} else {
				return true
			}
		}
	}
	return false
}
