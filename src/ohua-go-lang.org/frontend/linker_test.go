package frontend

import (
	"fmt"
	"go/importer"
	"go/types"
	"reflect"
	"testing"
)

// how would one like to write a stateful function in Go???

type State struct{
	some_value int
}

func (s State) Foo() int {
	fmt.Println("Running Foo ...")
	// access the state here
	s.some_value = 5
	return 5 + 1
}

func Some_independent_function(){

}

func Link(p string, fn_name string){
	fmt.Printf("Trying to link package: %s\n", p)
	pkg, err := importer.Default().Import(p)
	if err != nil {
		fmt.Printf("error: %s\n", err.Error())
		return
	}
	for _, declName := range pkg.Scope().Names() {
		fmt.Println(declName)
	}

	// THERE IS NO SUCH THING AS Class.forName in Go!
	// https://groups.google.com/forum/#!topic/golang-nuts/kTyvvFe8Bd8

	// but we find the (stateful) function
	var t types.Object
	t = pkg.Scope().Lookup("State")
	obj, _, indirect := types.LookupFieldOrMethod(t.Type(), true, pkg, "Foo") // check for a function in a type
	fmt.Printf("indirect: %b , obj: %s", indirect, obj)

	// TODO Create code to instantiate this struct and call the function on it.

	t.Type()
	fmt.Printf("Found struct: %s\n", t.Type())
	var rt reflect.Type
	rt = reflect.TypeOf(t.Type())
	fmt.Printf("Found type: %s\n", rt)
	//x := reflect.New(t.Type())
	//fmt.Println(x.Type())

	// how it should be:
	var a State
	k := reflect.TypeOf(a)
	fmt.Println("Found type: %+v", k)
	z := reflect.New(k)
	fmt.Println(z.Type())

	//xn := State(t)
	//fmt.Println("Found instance: %+v", xn)
	//xn.Foo()
	//y := reflect.ValueOf(x)
	//m := y.MethodByName("Foo")
	//fmt.Println(reflect.TypeOf(m))
	//for i := 0; i < y.NumMethod(); i++ {
	//	method := y.Method(i)
	//	fmt.Println(method.Name)
	//	//fmt.Println(method)
	//}
	//x := reflect.New(rt)
	//fmt.Println("Found value: %+v", x)
	//fmt.Println(rt.Name())
	//m,b := rt.MethodByName("Foo")
	//fmt.Println("Bool: %b", b)
	//fmt.Println("Found method: %+v", m)
	//
	//var i []reflect.Value
	////i[0] = reflect.ValueOf(5)
	//fmt.Println("Executing the call ...")
	//r := m.Call(i)
	//fmt.Println("Result: %i", r)
}

func Some_other_function(){

}

func Test_linking(t *testing.T){
	var b bool
	// FIXME I do not have a clue yet how these test files are being linked. I do not understand how one can create test
	//       data then.
	b = Exists("ohua-go-lang.org/frontend", "Foo")
	fmt.Printf("result: %b", b)
	if b == false { t.Fail() }
}