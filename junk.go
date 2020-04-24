package main

import (
	"fmt"

	"github.com/undeconstructed/gooo/cmps"
	"github.com/undeconstructed/gooo/tags"
	"github.com/undeconstructed/gooo/web"
)

func init() {
	web.Define(bodyComponent{}, buttonComponent{}, holderComponent{}, thingComponent{})
}

type body struct {
	web.Definition
}

type bodyComponent struct {
	web.State
	Args    body
	counter int
}

func (c *bodyComponent) Mount() {
	fmt.Printf("body mounting\n")
}

func (c *bodyComponent) Render() web.Output {
	fmt.Printf("body render: %v\n", c.Args)
	return tags.Body().With(
		tags.H1().With(
			web.Text("Title"),
		),
		cmps.Label{Text: fmt.Sprintf("foo count: %d", c.counter)},
		tags.P().With(
			web.Textf("count: %d", c.counter),
		),
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
	c.counter++
	c.Update()
}

type holder struct {
	web.Definition
	Child web.Definition
}

type holderComponent struct {
	web.Component
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
	web.Component
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
	web.Component
	Args button
}

func (c *buttonComponent) Render() web.Output {
	return web.Tag("button").With(
		web.Text(c.Args.Label),
	).On("click", c.Args.OnClick)
}
