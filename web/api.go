package web

import "fmt"

// Any thing.
type Any interface{}

// Output is rendered stuff.
type Output interface {
	isOutput()
}

// text struct is static text to render
type text struct {
	s string
}

// mark *text as Output type
func (*text) isOutput() {
}

// Text node
func Text(s string) Output {
	return &text{s: s}
}

// Textf node
func Textf(format string, a ...interface{}) Output {
	return &text{s: fmt.Sprintf(format, a...)}
}

// HTML is HTML output.
type HTML struct {
	tag      string
	events   []eventHandler
	children []Output
}

// Tag tag
func Tag(tag string) *HTML {
	return &HTML{tag: tag}
}

// With children
func (h *HTML) With(o ...Output) *HTML {
	h.children = append(h.children, o...)
	return h
}

// On defines an event handler
func (h *HTML) On(event string, handler EventHandler) *HTML {
	h.events = append(h.events, eventHandler{event, handler})
	return h
}

// mark *HTML as Output type
func (*HTML) isOutput() {
}

// internalComponent is an internal marker for any type of component
type internalComponent interface {
	isComponent()
}

// Component marks a component. Just embed one.
type Component struct {
}

// implement Component.
func (*Component) isComponent() {
}

// Renderer can render itelf, implement it on components.
type Renderer interface {
	Render() Output
}

// Mounter does something on mount..
type Mounter interface {
	Mount()
}

// Unmounter does something on unmount.
type Unmounter interface {
	Unmounter()
}

// Stateful is for stateful components. Embed a State to implement it.
type Stateful interface {
	set(s State)
	Update()
}

// State lets a component be stateful.
type State struct {
	Component
	render func()
}

// implement Stateful.
func (s *State) set(n State) {
	*s = n
}

// Update is for component to have itself re-rendered.
func (s *State) Update() {
	// fmt.Println("update")
	s.render()
}

// Definition is for component def typeRegistry to self define.
type Definition interface {
	isOutput()
}

// EventHandler is
type EventHandler func()

// internal view of eventHandler
type eventHandler struct {
	event   string
	handler EventHandler
}
