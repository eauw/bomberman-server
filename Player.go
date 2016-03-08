package main

// Player
type Player struct {
	id     string
	name   string
	points int
}

// NewPlayer function is the players constructor
func NewPlayer(name string) *Player {
	playerID := randomString(8)

	return &Player{
		id:   playerID,
		name: name,
	}
}

func (player *Player) setName(name string) {
	player.name = name
}
