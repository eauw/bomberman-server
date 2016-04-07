package gamemanager

import (
	"bomberman-server/helper"
	"fmt"
	"log"
)

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
	fields := gameMap.GetNOSWFieldsOfField(bomb.field)

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
			}

			p.resetSpecials()
		}
	}

	// Bombe nach Explosion wieder verfügbar machen

	bomb.timer = 5
	bomb.field.bombs = []*Bomb{}
	bomb.field = nil
	log.Println(fmt.Sprintf("%s", gameMap.bombs))
	gameMap.removeBomb(bomb)
	log.Println(fmt.Sprintf("%s", gameMap.bombs))

	for _, f := range fields {
		for _, b := range f.bombs {
			b.explode(gameMap)
		}
	}
}
