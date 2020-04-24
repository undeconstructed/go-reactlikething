package tags

import (
	"github.com/undeconstructed/gooo/web"
)

// H1 is
func H1() *web.HTML {
	return web.Tag("h1")
}

// P is
func P() *web.HTML {
	return web.Tag("p")
}

// Body is
func Body() *web.HTML {
	return web.Tag("body")
}

// Div is
func Div() *web.HTML {
	return web.Tag("div")
}

// Button is
func Button() *web.HTML {
	return web.Tag("button")
}

// OL is
func OL() *web.HTML {
	return web.Tag("ol")
}
