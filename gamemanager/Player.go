package gamemanager

import (
	"bomberman-server/helper"
	"fmt"
)

// Player
type Player struct {
	game         *Game
	id           string
	ip           string
	name         string
	points       int
	currentField *Field
	isParalyzed  int
	isFox        int
	bombs        []*Bomb
	throwrange   int
	protection   int // Dauer des Schutzes falls meinen welchen eingesammelt hat
}

// NewPlayer function is the players constructor
func NewPlayer(n string, f *Field) *Player {
	playerID := helper.RandomString(8)

	return &Player{
		id:           playerID,
		name:         n,
		currentField: f,
		isParalyzed:  0,
		isFox:        0,
		bombs:        make([]*Bomb, 0),
		throwrange:   1,
		protection:   0,
	}
}

func (player *Player) GetIP() string {
	return player.ip
}

func (player *Player) SetIP(ip string) {
	player.ip = ip
}

func (player *Player) GetID() string {
	return player.id
}

func (player *Player) SetID(id string) {
	player.id = id
}

func (player *Player) GetName() string {
	return player.name
}

func (player *Player) SetName(name string) {
	player.name = name
}

func (player *Player) GetCurrentField() *Field {
	return player.currentField
}

func (player *Player) SetCurrentField(field *Field) {
	player.currentField = field
}

func (player *Player) applySpecial(special *Special) {
	switch special.powerType {
	case "r":
		player.throwrange += 1
		break

	case "b":
		player.addBomb()
		break

	case "h":
		player.protection = 5
		break
	}
}

func (player *Player) addBomb() {
	newBomb := NewBomb()
	newBomb.owner = player
	player.bombs = append(player.bombs, newBomb)
}

func (player *Player) getAvailableBomb() *Bomb {
	for i := range player.bombs {
		if player.bombs[i].isPlaced == false {
			return player.bombs[i]
		}
	}

	return nil
}

func (player *Player) resetSpecials() {
	player.throwrange = 1
	player.protection = 0
	player.bombs = []*Bomb{}
	player.addBomb()
}

func (player *Player) toString() string {
	idString := ""
	ipString := ""
	nameString := ""
	//pointsString := ""
	currentFieldString := ""

	if player.id != "" {
		idString = player.id
	} else {
		idString = "nil"
	}

	if player.ip != "" {
		ipString = player.ip
	} else {
		ipString = "nil"
	}

	if player.name != "" {
		nameString = player.name
	} else {
		nameString = "nil"
	}

	if player.currentField != nil {
		currentFieldString = player.currentField.toString()
	} else {
		currentFieldString = "nil"
	}

	playerString := fmt.Sprintf("id: " + idString + " ip: " + ipString + " name: " + nameString + " currentField: " + currentFieldString)

	return playerString
}
