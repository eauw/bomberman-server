package main

// Player
type Player struct {
	name   string
	points int
}

// NewPlayer function is the players constructor
func NewPlayer(name string) *Player {
	return &Player{
		name: name,
	}
}

func (player *Player) setName(name string) {
	player.name = name
}
