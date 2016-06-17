package gamemanager

import (
	"fmt"

	"github.com/eauw/bomberman-server/helper"
	"github.com/fatih/color"
)

// Player ...
type Player struct {
	id           string
	ip           string
	name         string
	points       int
	currentField *Field
	isParalyzed  int
	isFox        bool
	foxRounds    int
	bombs        []*Bomb
	reach        int
	throwrange   int
	protection   int    // Dauer des Schutzes falls man welchen eingesammelt hat
	msg          string // Nachricht, die dem Spieler nach dem naechsten gamestate geschickt werden soll.
	color        *color.Color
}

// NewPlayer function is the players constructor
func NewPlayer(name string) *Player {
	playerID := helper.GeneratePlayerID()

	return &Player{
		id:          playerID,
		name:        name,
		isParalyzed: 0,
		isFox:       false,
		foxRounds:   0,
		bombs:       make([]*Bomb, 0),
		reach:       1,
		throwrange:  9,
		protection:  0,
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
		player.reach += 1

	case "b":
		player.addBomb()

	case "h":
		player.protection = 20
	}
}

func (player *Player) addBomb() {
	newBomb := NewBomb()
	newBomb.owner = player
	player.bombs = append(player.bombs, newBomb)
}

func (player *Player) getAvailableBomb() *Bomb {
	for _, b := range player.bombs {
		if b.field == nil {
			return b
		}
	}

	return nil
}

func (player *Player) resetSpecials() {
	player.reach = 1
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
