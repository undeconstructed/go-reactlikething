package main

import (
	"fmt"

	"github.com/undeconstructed/gooo/tags"
	"github.com/undeconstructed/gooo/web"
)

type body struct {
	web.Definition
}

type bodyComponent struct {
	Args body
}

func (c *bodyComponent) Render() web.Output {
	return tags.Div().With(
		thing{Foo: "testing1"},
		tags.Div().With(
			thing{Foo: "testing2"},
		),
		thing{Foo: "testing3"},
		holder{Child: thing{Foo: "???"}},
		tags.P().With(
			button{
				Label:   "click me",
				OnClick: c.onButtonClick,
			},
		),
	)
}

func (c *bodyComponent) onButtonClick() {
	fmt.Println("clicky!")
}

type holder struct {
	web.Definition
	Child web.Definition
}

type holderComponent struct {
	Args holder
}

func (h *holderComponent) Render() web.Output {
	return tags.Div().With(
		h.Args.Child,
	)
}

type thing struct {
	web.Definition
	Foo string
}

type thingComponent struct {
	Args thing
}

func (c *thingComponent) Render() web.Output {
	return tags.P().With(
		web.Text(c.Args.Foo),
	)
}

type button struct {
	web.Definition
	Label   string
	OnClick web.EventHandler
}

type buttonComponent struct {
	Args button
}

func (c *buttonComponent) Render() web.Output {
	return web.Tag("button").With(
		web.Text(c.Args.Label),
	).On("click", c.Args.OnClick)
}

func main() {
	sys := web.New()
	sys.Define(bodyComponent{}, buttonComponent{}, holderComponent{}, thingComponent{})

	sys.MainBody(body{})
}
