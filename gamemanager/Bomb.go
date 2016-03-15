package gamemanager

type Bomb struct {
	owner *Player
	field *Field
}

func NewBomb(p *Player, f *Field) *Bomb {
	return &Bomb{
		owner: p,
		field: f,
	}
}

func (bomb *Bomb) explode(gameMap *GameMap) {
	fields := gameMap.GetNOSWFieldsOfField(bomb.field)

	// Ausgangsfeld hinzufügen
	fields = append(fields, bomb.field)

	for i := range fields {
		// Wände werden durch Explosionsstrahl zerstört
		fields[i].containsWall = false

		// Specials werden durch Explosionsstrahl zerstört
		fields[i].containsSpecial = false

		// Spieler werden durch Explosionsstrahl gelähmt
		for pI := range fields[i].players {
			fields[i].players[pI].isParalyzed = true
		}
	}
}
