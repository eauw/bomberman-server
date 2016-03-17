package gamemanager

type Bomb struct {
	owner *Player
	field *Field // das Feld auf das die Bombe geworfen wurde
	reach int    // Reichweite der Bombe
	timer int    // Zeit/Runden die die Bombe braucht zum explodieren
}

func NewBomb() *Bomb {
	return &Bomb{
		reach: 1,
		timer: 3,
	}
}

func (bomb *Bomb) explode(gameMap *GameMap) {
	fields := gameMap.GetNOSWFieldsOfField(bomb.field)

	// Ausgangsfeld hinzufügen
	fields = append(fields, bomb.field)

	for i := range fields {
		// Wände werden durch Explosionsstrahl zerstört
		if fields[i].wall != nil {
			if fields[i].wall.isDestructible {
				fields[i].wall = nil
			}
		}

		// Specials werden durch Explosionsstrahl zerstört
		fields[i].wall = nil

		// Spieler werden durch Explosionsstrahl gelähmt
		for _, v := range fields[i].players {
			v.isParalyzed = true
		}
	}
}
