package cmps

import (
	"fmt"

	"github.com/undeconstructed/gooo/tags"
	"github.com/undeconstructed/gooo/web"
)

type Label struct {
	web.Definition
	Text string
}

type LabelComponent struct {
	web.Component
	Args Label
}

func (sc *LabelComponent) Render() web.Output {
	fmt.Printf("label render: %v\n", sc.Args)
	return tags.P().With(web.Text(sc.Args.Text))
}

func init() {
	web.Define(LabelComponent{})
}
