package web

type Static struct {
	Definition
	Out Output
}

type StaticComponent struct {
	Component
	Args Static
}

func (sc *StaticComponent) Render() Output {
	return sc.Args.Out
}

func init() {
	Define(StaticComponent{})
}
