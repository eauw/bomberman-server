package gamemanager

type Round struct {
	id             int
	playerCommands map[string]string
}

func NewRound() *Round {
	return &Round{
		playerCommands: map[string]string{},
	}
}
