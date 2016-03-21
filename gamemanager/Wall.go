package gamemanager

type Wall struct {
	isDestructible bool
}

func NewWall(destructible bool) *Wall {
	return &Wall{
		isDestructible: destructible,
	}
}
