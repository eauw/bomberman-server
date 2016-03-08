package main

type Position struct {
	x int
	y int
}

func NewPosition() *Position {
	return &Position{
		x: 0,
		y: 0,
	}
}

func (p *Position) setPosition(x int, y int) {
	p.x = x
	p.y = y
}
