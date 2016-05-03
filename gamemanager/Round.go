package gamemanager

type Round struct {
	id             int
	playerCommands map[string]PlayerCommand
}

func NewRound() *Round {
	return &Round{
		playerCommands: map[string]PlayerCommand{},
	}
}
