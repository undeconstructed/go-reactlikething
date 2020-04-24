package main

import (
	"github.com/undeconstructed/gooo/tags"
	"github.com/undeconstructed/gooo/web"
)

func init() {
	web.Define(squareComponent{}, boardComponent{}, gameComponent{})
}

type square struct {
	web.Definition
	Value   string
	OnClick web.EventHandler
}

type squareComponent struct {
	web.State
	Args square
}

func (c *squareComponent) Render() web.Output {
	return tags.Button().Class("square").With(
		web.Textf("%s", c.Args.Value),
	).On("click", c.Args.OnClick)
}

// func (c *squareComponent) onClick() {
// 	c.Update()
// }

type board struct {
	web.Definition
}

type boardComponent struct {
	web.State
	Args    board
	Values  []string
	XIsNext bool
}

func (c *boardComponent) Init() {
	c.Values = make([]string, 9, 9)
	c.XIsNext = true
}

func (c *boardComponent) renderSquare(i int) web.Output {
	return square{Value: c.Values[i], OnClick: func() {
		c.onSquareClick(i)
	}}
}

func (c *boardComponent) onSquareClick(i int) {
	if calculateWinner(c.Values) != "" || c.Values[i] != "" {
		return
	}

	if c.XIsNext {
		c.Values[i] = "X"
		c.XIsNext = false
	} else {
		c.Values[i] = "O"
		c.XIsNext = true
	}

	c.Update()
}

func (c *boardComponent) Render() web.Output {
	winner := calculateWinner(c.Values)
	status := "Next player: X"
	if winner != "" {
		status = "Winner: " + winner
	} else if !c.XIsNext {
		status = "Next player: O"
	}

	return tags.Div().With(
		tags.Div().Class("status").With(
			web.Text(status),
		),
		tags.Div().Class("board-row").With(
			c.renderSquare(0),
			c.renderSquare(1),
			c.renderSquare(2),
		),
		tags.Div().Class("board-row").With(
			c.renderSquare(3),
			c.renderSquare(4),
			c.renderSquare(5),
		),
		tags.Div().Class("board-row").With(
			c.renderSquare(6),
			c.renderSquare(7),
			c.renderSquare(8),
		),
	)
}

var lines = [][]int{
	[]int{0, 1, 2},
	[]int{3, 4, 5},
	[]int{6, 7, 8},
	[]int{0, 3, 6},
	[]int{1, 4, 7},
	[]int{2, 5, 8},
	[]int{0, 4, 8},
	[]int{2, 4, 6},
}

func calculateWinner(squares []string) string {
	for _, line := range lines {
		a, b, c := line[0], line[1], line[2]
		if squares[a] != "" && squares[a] == squares[b] && squares[a] == squares[c] {
			return squares[a]
		}
	}
	return ""
}

type game struct {
	web.Definition
}

type gameComponent struct {
	web.State
	Args game
}

func main() {
	web.AddStyleSheet("style.css")

	web.MainBody(tags.Body().With(
		tags.Div().Class("game").With(
			tags.Div().Class("game-board").With(
				board{},
			),
			tags.Div().Class("game-info").With(
				tags.Div(),
				tags.OL(),
			),
		),
	))
}
