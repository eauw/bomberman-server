package gamemanager

type Bomb struct {
	owner *Player
}

func NewBomb(player *Player) *Bomb {
	return &Bomb{
		owner: player,
	}
}

func (bomb *Bomb) explode() {

}
