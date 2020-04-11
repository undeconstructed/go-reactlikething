package web

import (
	"fmt"
	"reflect"
	"syscall/js"
)

// Any thing.
type Any interface{}

// Output is rendered stuff.
type Output interface {
	isOutput()
}

type text struct {
	s string
}

func (*text) isOutput() {
}

// Text node
func Text(s string) Output {
	return &text{s: s}
}

// HTML is HTML output.
type HTML struct {
	tag      string
	events   map[string]EventHandler
	children []Output
}

// Tag tag
func Tag(tag string) *HTML {
	return &HTML{tag: tag}
}

// With adds children
func (h *HTML) With(o ...Output) *HTML {
	h.children = append(h.children, o...)
	return h
}

// On defines an event handler
func (h *HTML) On(event string, handler EventHandler) *HTML {
	if h.events == nil {
		h.events = map[string]EventHandler{}
	}
	h.events[event] = handler
	return h
}

func (*HTML) isOutput() {
}

// Component can render itelf, and will persist in the node tree until it is unwanted.
type Component interface {
	Render() Output
}

// Definition is for component def types to self define.
type Definition interface {
	isOutput()
}

// EventHandler is
type EventHandler func()

type types map[reflect.Type]reflect.Type

func (ts types) add(a, c Any) {
	ts[reflect.TypeOf(a)] = reflect.TypeOf(c)
}

func (ts types) inst(a Any) Component {
	ct := ts[reflect.TypeOf(a)]
	cv := reflect.New(ct)
	reflect.Indirect(cv).FieldByName("Args").Set(reflect.ValueOf(a))
	return cv.Interface().(Component)
}

type node struct {
	tag      string
	text     string
	cmp      Component
	events   map[string]EventHandler
	children []*node
	next     *node
}

func indent(depth int) {
	for i := 0; i < depth; i++ {
		fmt.Print("  ")
	}
}

func (n *node) printout(depth int) {
	indent(depth)
	fmt.Printf("tag: %s, text: %s, cmp: %v\n", n.tag, n.text, n.cmp)
	for _, c := range n.children {
		c.printout(depth + 1)
	}
}

// System is core.
type System struct {
	types types
	root  *node
}

// New news.
func New() *System {
	return &System{
		types: types{},
	}
}

// Define links args to components.
func (s *System) Define(cmps ...Any) {
	for _, cmp := range cmps {
		ct := reflect.TypeOf(cmp)
		if ct.Kind() != reflect.Struct {
			panic("cmp must be struct: " + ct.Name())
		}
		tic := reflect.TypeOf((*Component)(nil)).Elem()
		if !reflect.PtrTo(ct).Implements(tic) {
			panic("*cmp must implement Component: " + ct.Name())
		}
		argsField, exists := ct.FieldByName("Args")
		if !exists {
			panic("cmp must have field Args: " + ct.Name())
		}
		dt := argsField.Type
		if dt.Kind() != reflect.Struct {
			panic("def must be struct: " + ct.Name())
		}
		// for i := 0; i < dt.NumMethod(); i++ {
		// 	f := dt.Method(i)
		// 	fmt.Printf("m %v\n", f)
		// }
		// for i := 0; i < dt.NumField(); i++ {
		// 	f := dt.Field(i)
		// 	fmt.Printf("f %v\n", f)
		// }
		tid := reflect.TypeOf((*Definition)(nil)).Elem()
		if !dt.Implements(tid) {
			panic("def must have Definition: " + ct.Name())
		}
		// _, exists = dt.FieldByName("Definition")
		// if !exists {
		// 	panic("def must have Definition: " + ct.Name())
		// }
		s.types[dt] = ct
	}
}

// MainBody sets component root and then lives forever.
func (s *System) MainBody(body Output) {
	// make empty root
	root := &node{}
	// now reconcile and render through the tree
	s.renderToNode(body, root)
	// now put into the page
	p := js.Global().Get("document").Get("body")
	s.renderToPage(root, p)

	// root.printout(0)

	select {}
}

func (s *System) renderToNode(out Output, n *node) {
	switch v := out.(type) {
	case *text:
		n.text = v.s
	case *HTML:
		n.tag = v.tag
		for _, c := range v.children {
			n2 := &node{}
			n.events = v.events
			n.children = append(n.children, n2)
			s.renderToNode(c, n2)
		}
	case Definition:
		// TODO = reuse component
		c := s.types.inst(v)
		n.cmp = c
		out2 := c.Render()
		s.renderToNode(out2, n)
	default:
		panic("fail")
	}
}

func (s *System) renderToPage(n *node, target js.Value) {
	d := js.Global().Get("document")
	makeEl := func(tag string) js.Value {
		return d.Call("createElement", tag)
	}
	var inner func(n *node, p js.Value)
	inner = func(n *node, p js.Value) {
		if n.tag == "" {
			if n.text != "" {
				p.Set("textContent", n.text)
			}
			return
		}
		e := makeEl(n.tag)
		for ev, h := range n.events {
			fmt.Printf("ev %s h %v\n", ev, h)
			wrapper := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
				h()
				return nil
			})
			e.Call("addEventListener", ev, wrapper)
		}
		for _, c := range n.children {
			inner(c, e)
		}
		p.Call("append", e)
	}
	inner(n, target)
}
