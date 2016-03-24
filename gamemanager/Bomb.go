package gamemanager

type Bomb struct {
	owner    *Player
	field    *Field // das Feld auf das die Bombe geworfen wurde
	reach    int    // Reichweite der Bombe
	timer    int    // Zeit/Runden die die Bombe braucht zum explodieren
	isPlaced bool
}

func NewBomb() *Bomb {
	return &Bomb{
		reach:    1,
		timer:    3,
		isPlaced: false,
	}
}

func (bomb *Bomb) explode(gameMap *GameMap) {
	fields := gameMap.GetNOSWFieldsOfField(bomb.field)

	// Ausgangsfeld hinzufügen
	fields = append(fields, bomb.field)

	for i := range fields {
		// Feld als "gerade am explodieren" markieren
		fields[i].explodes = true

		// Wände werden durch Explosionsstrahl zerstört
		if fields[i].wall != nil {
			if fields[i].wall.isDestructible {
				fields[i].wall = nil
				fields[i].explodes = true
			} else {
				fields[i].explodes = false
			}
		}

		// Specials werden durch Explosionsstrahl zerstört
		fields[i].special = nil

		// Spieler werden durch Explosionsstrahl gelähmt
		for _, v := range fields[i].players {
			if v.protection == false {
				v.isParalyzed = true
			}

			v.resetSpecials()
		}
	}

	// Bombe nach Explosion wieder verfügbar machen
	bomb.isPlaced = false
	bomb.field.bombs = []*Bomb{}
	gameMap.removeBomb(bomb)

	for _, f := range fields {
		if len(f.bombs) > 0 {
			for _, b := range f.bombs {
				b.explode(gameMap)
			}
		}
	}
}
