package gamemanager

import (
	"github.com/eauw/bomberman-server/helper"
)

const bombtimer = 5

type Bomb struct {
	owner *Player
	field *Field // das Feld auf das die Bombe geworfen wurde
	reach int    // Reichweite der Bombe
	timer int    // Zeit/Runden die die Bombe braucht zum explodieren
	id    string
}

func NewBomb() *Bomb {
	id := helper.GenerateBombID()

	return &Bomb{
		reach: 1,
		timer: 5,
		id:    id,
		field: nil,
		owner: nil,
	}
}

func (bomb *Bomb) explode(gameMap *GameMap) {
	fields := gameMap.GetFieldsThatAreImpactedFromExplosion(bomb.field, bomb.owner.reach)

	// Ausgangsfeld hinzufügen da die Methode GetNOSWFieldsOfField() nur die Felder links, oben, rechts und unten holt
	fields = append(fields, bomb.field)

	for _, f := range fields {
		// Feld als "gerade am explodieren" markieren
		f.explodes = true

		// Wände werden durch Explosionsstrahl zerstört
		if f.wall != nil {
			if f.wall.isDestructible {
				f.wall = nil
				f.explodes = true
			} else {
				f.explodes = false
			}
		}

		// Specials werden durch Explosionsstrahl zerstört
		f.special = nil

		// Spieler werden durch Explosionsstrahl gelähmt
		for _, p := range f.players {
			if p.protection == 0 {
				p.isParalyzed = 3
				p.resetSpecials()
			}
		}
	}

	// Bombe nach Explosion wieder verfügbar machen

	bomb.timer = bombtimer
	bomb.field.bombs = []*Bomb{}
	bomb.field = nil
	gameMap.removeBomb(bomb)

	for _, f := range fields {
		for _, b := range f.bombs {
			b.explode(gameMap)
		}
	}
}
