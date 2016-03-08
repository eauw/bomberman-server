package main

import "bomberman-server/helper"

// Player
type Player struct {
	id           string
	name         string
	points       int
	currentField *Field
}

// NewPlayer function is the players constructor
func NewPlayer(name string) *Player {
	playerID := helper.RandomString(8)

	return &Player{
		id:   playerID,
		name: name,
	}
}

func (player *Player) setName(name string) {
	player.name = name
}

// func (player *Player) setPosition(x int, y int) {
// 	player.position.setPosition(x, y)
// }
