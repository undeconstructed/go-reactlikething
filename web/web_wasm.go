// +build wasm

package web

import (
	"fmt"
	"reflect"
	"syscall/js"
)

// helpers

var document = js.Global().Get("document")

func createElement(tag string) domElement {
	n := document.Call("createElement", tag)
	// fmt.Printf("created %s %#v\n", tag, n)
	return domElement(n)
}

func createTextNode(text string) domElement {
	n := document.Call("createTextNode", text)
	return domElement(n)
}

type domElement js.Value

func (de domElement) IsUndefined() bool {
	return js.Value(de).IsUndefined()
}

func (de domElement) Append(c domElement) {
	// fmt.Printf("%#v append %#v\n", de, c)
	js.Value(de).Call("append", js.Value(c))
}

func (de domElement) RemoveChild(c domElement) {
	// fmt.Printf("%#v removeChild %#v\n", de, c)
	js.Value(de).Call("removeChild", js.Value(c))
}

func (de domElement) ReplaceChild(c1 domElement, c0 domElement) {
	js.Value(de).Call("replaceChild", js.Value(c1), js.Value(c0))
}

func (de domElement) AddClass(c string) {
	js.Value(de).Get("classList").Call("add", c)
}

func (de domElement) Listen(e string, f js.Func) {
	js.Value(de).Call("addEventListener", e, f)
}

// nodes

type Node interface {
	// Accept responds whether a node would accept a new definition
	Accept(Output) bool
	// Update accepts a new definition, but doesn't act on it
	Update(Output)
	// Layout makes DOM changes
	Layout(domElement)
	// Unlayout undoes DOM changes
	Unlayout(domElement)
}

func makeNode(out Output) Node {
	switch out := out.(type) {
	case *HTML:
		return &tagNode{def: out}
	case *text:
		return &textNode{def1: out}
	case Definition:
		return &componentNode{def1: out}
	default:
		panic("unknown thing")
	}
}

type componentNode struct {
	// ctype    componentType
	def0     Definition
	def1     Definition
	cmp      internalComponent
	child    Node
	oldChild Node
}

func (n *componentNode) Accept(out Output) bool {
	if reflect.TypeOf(out) == reflect.TypeOf(n.def0) {
		// accept a new definition of the same type
		return true
	}
	return false
}

func (n *componentNode) Update(out Output) {
	def1 := out.(Definition)
	// maybe not comparable
	// if def1 == n.def0 {
	// 	return
	// }
	n.def1 = def1
}

func (n *componentNode) Render() []*componentNode {
	// NB def1 is already swapped into def0, and cmp updated

	// if this component actually renders ...
	if re, ok := n.cmp.(Renderer); ok {
		out := re.Render()
		if n.child == nil {
			// shortcut first render
			n.child = makeNode(out)
		} else {
			if n.child.Accept(out) {
				n.child.Update(out)
			} else {
				if n.oldChild == nil {
					// only keep the oldest child, any other hasn't been laid out ever
					n.oldChild = n.child
				}
				n.child = makeNode(out)
			}
		}
	}

	var childComponents []*componentNode

	// if we have a child now
	if n.child != nil {
		switch c := n.child.(type) {
		case *componentNode:
			childComponents = append(childComponents, c)
		case *tagNode:
			moreChildComponents := c.Render()
			// TODO - maybe map in existing components
			childComponents = append(childComponents, moreChildComponents...)
		case *textNode:
			c.Render()
		}
	}

	// TODO - match up any new child nodes with abandoned ones

	return childComponents
}

func (n *componentNode) Layout(parent domElement) {
	if n.oldChild != nil {
		n.oldChild.Unlayout(parent)
		n.oldChild = nil
	}

	if n.child != nil {
		n.child.Layout(parent)
	}
}

func (n *componentNode) Unlayout(parent domElement) {
	if n.child != nil {
		n.child.Unlayout(parent)
	}
}

type tagNode struct {
	def      *HTML
	children []Node
	jso      domElement
}

func (n *tagNode) Accept(out Output) bool {
	if h, ok := out.(*HTML); ok {
		if h.tag == n.def.tag {
			return true
		}
	}
	return false
}

func (n *tagNode) Update(out Output) {
	h := out.(*HTML)
	n.def = h
}

func (n *tagNode) Render() []*componentNode {
	var childComponents []*componentNode

	recurse := func(child Node) {
		switch v := child.(type) {
		case *componentNode:
			childComponents = append(childComponents, v)
		case *tagNode:
			moreChildComponents := v.Render()
			childComponents = append(childComponents, moreChildComponents...)
		case *textNode:
			v.Render()
		}
	}

	if n.children == nil {
		// shortcut first render
		for _, c := range n.def.children {
			child := makeNode(c)
			recurse(child)
			n.children = append(n.children, child)
		}
	} else {
		// TODO new children, gone children, moves, etc
		for i, child := range n.children {
			newDef := n.def.children[i]
			if child.Accept(newDef) {
				child.Update(newDef)
				recurse(child)
			} else {
				// this part of the tree is broken
				panic("failed to update html tree")
			}
		}
	}

	// TODO return abandoned children as well
	return childComponents
}

func (n *tagNode) Layout(parent domElement) {
	if n.jso.IsUndefined() {
		n.jso = createElement(n.def.tag)
		for _, evh := range n.def.events {
			e := evh.event
			h := evh.handler
			wrapper := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
				h()
				return nil
			})
			n.jso.Listen(e, wrapper)
		}
		for _, c := range n.def.classes {
			n.jso.AddClass(c)
		}
		parent.Append(n.jso)
	}

	// TODO - update tag
	// TODO - remove old children

	for _, c := range n.children {
		c.Layout(n.jso)
	}
}

func (n *tagNode) Unlayout(parent domElement) {
	for _, c := range n.children {
		c.Unlayout(n.jso)
	}
	// TODO remove listeners
	parent.RemoveChild(n.jso)
}

type textNode struct {
	def0 *text
	def1 *text
	jso  domElement
}

func (n *textNode) Accept(out Output) bool {
	if _, ok := out.(*text); ok {
		return true
	}
	return false
}

func (n *textNode) Update(out Output) {
	def1 := out.(*text)
	if def1 == n.def0 || def1.s == n.def0.s {
		return
	}
	n.def1 = def1
}

func (n *textNode) Render() {
	// nothing to do
}

func (n *textNode) Layout(parent domElement) {
	if n.jso.IsUndefined() {
		// if never laid out before
		n.def0 = n.def1
		n.def1 = nil
		n.jso = createTextNode(n.def0.s)
		parent.Append(n.jso)
		return
	} else if n.def1 != nil {
		n.def0 = n.def1
		n.def1 = nil
		oldJso := n.jso
		n.jso = createTextNode(n.def0.s)
		parent.ReplaceChild(n.jso, oldJso)
	}
}

func (n *textNode) Unlayout(parent domElement) {
	if !n.jso.IsUndefined() {
		parent.RemoveChild(n.jso)
	}
}

// system

type sys struct {
	renderCh chan *componentNode
	mountCh  chan *componentNode
}

func (s *sys) render(node *componentNode) {
	s.renderCh <- node
}

func (s *sys) mount(node *componentNode) {
	s.mountCh <- node
}

func addStyleSheet(url string) {
	css := js.Value(createElement("link"))
	css.Set("rel", "stylesheet")
	css.Set("type", "text/css")
	css.Set("href", url)
	document.Get("head").Call("append", css)
}

func mainBody(def Definition) {
	renderCh := make(chan *componentNode, 100)
	mountCh := make(chan *componentNode, 100)
	s := &sys{
		renderCh: renderCh,
		mountCh:  mountCh,
	}

	root := rootNode(def)

	go func() {
		for {
			fmt.Printf("render wait\n")

			// wait under something needs rendering
			n := <-renderCh
			// render first thing in the queue
			renderComponent(s, n)

			// now loop until render queue is empty
		renderQ:
			for {
				select {
				case n := <-renderCh:
					renderComponent(s, n)
				default:
					break renderQ
				}
			}

			// TODO - unmount anything old

			// and layout the whole lot
			root.Layout(domElement{})

			// mount anything new
		mountQ:
			for {
				select {
				case n := <-mountCh:
					// logNode("mount?", n)
					if c, ok := n.cmp.(Mounter); ok {
						c.Mount()
					}
				default:
					break mountQ
				}
			}
		}
	}()

	s.render(root)
}

func renderComponent(s *sys, n *componentNode) {
	if n.cmp == nil {
		ct, _ := types.get(n.def1)
		rerender := func() {
			// fmt.Printf("rerender\n")
			s.render(n)
		}
		n.cmp = ct.new(State{render: rerender})

		if ir, ok := n.cmp.(Initer); ok {
			ir.Init()
		}

		// to be mounted later, after layout performed
		s.mount(n)
	}

	if n.def1 != nil {
		ct, _ := types.get(n.def1)
		ct.update(n.cmp, n.def1)
		n.def0 = n.def1
		n.def1 = nil
		// TODO - inform of this change?
	}

	childComponents := n.Render()

	for _, c := range childComponents {
		s.render(c)
	}
}

func rootNode(def Definition) *componentNode {
	// html always has an implicit body, so fetch it
	body := document.Get("body")
	// fmt.Printf("body %#v\n", body)
	// make a node to hold it
	bodyNode := &tagNode{
		def: Tag("body"),
		jso: domElement(body),
	}
	// make the root node with body as only child
	rootNode := &componentNode{
		def1:  def,
		child: bodyNode,
	}

	logNode("root", rootNode)

	return rootNode
}

func logNode(s string, n Node) {
	fmt.Printf("node: %[1]s %[2]p %#[2]v\n", s, n)
}
