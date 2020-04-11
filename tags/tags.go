package tags

import (
	"github.com/undeconstructed/gooo/web"
)

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
