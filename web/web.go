package web

import (
	"fmt"
	"reflect"
)

type componentType struct {
	aType reflect.Type
	cType reflect.Type
}

func (ct componentType) match(c internalComponent) bool {
	return reflect.TypeOf(c) == reflect.PtrTo(ct.cType)
}

func (ct componentType) new(s State) internalComponent {
	cv := reflect.New(ct.cType)
	newCmp := cv.Interface().(internalComponent)
	if st, ok := newCmp.(Stateful); ok {
		st.set(s)
	}
	return newCmp
}

func (ct componentType) update(c internalComponent, args Any) {
	cv := reflect.ValueOf(c)
	reflect.Indirect(cv).FieldByName("Args").Set(reflect.ValueOf(args))
}

type typeRegistry map[reflect.Type]componentType

func (ts typeRegistry) set(a, c reflect.Type) {
	ts[a] = componentType{
		aType: a,
		cType: c,
	}
}

func (ts typeRegistry) get(a Any) (componentType, bool) {
	ct, exists := ts[reflect.TypeOf(a)]
	return ct, exists
}

var types = typeRegistry{}

func indent(depth int) {
	for i := 0; i < depth; i++ {
		fmt.Print("  ")
	}
}

// Define links args to components.
func Define(cmps ...Any) {
	for _, cmp := range cmps {
		ct := reflect.TypeOf(cmp)
		if ct.Kind() != reflect.Struct {
			panic("cmp must be struct: " + ct.Name())
		}
		tic := reflect.TypeOf((*internalComponent)(nil)).Elem()
		if !reflect.PtrTo(ct).Implements(tic) {
			panic("*cmp must be a Component: " + ct.Name())
		}
		// tir := reflect.TypeOf((*Renderer)(nil)).Elem()
		// if !reflect.PtrTo(ct).Implements(tir) {
		// 	panic("*cmp must implement Renderer: " + ct.Name())
		// }
		argsField, exists := ct.FieldByName("Args")
		if !exists {
			panic("cmp must have field Args: " + ct.Name())
		}
		dt := argsField.Type
		if dt.Kind() != reflect.Struct {
			panic("def must be struct: " + ct.Name())
		}
		tid := reflect.TypeOf((*Definition)(nil)).Elem()
		if !dt.Implements(tid) {
			panic("def must be a Definition: " + ct.Name())
		}
		types.set(dt, ct)
	}
}

// MainBody renders the body - it be a body, or a def that always renders a body.
func MainBody(body Output) {
	var rootDef Definition

	switch v := body.(type) {
	case *HTML:
		if v.tag == "body" {
			rootDef = Static{Out: body}
		} else {
			panic("must render a body")
		}
	case Definition:
		rootDef = v
	}

	mainBody(rootDef)

	select {}
}
