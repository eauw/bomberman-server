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
	hasSpecial   bool // wenn true dann hat er ein special eingesammelt und kann währenddessen irgendwas besonderes
}

// NewPlayer function is the players constructor
func NewPlayer(name string) *Player {
	playerID := helper.RandomString(8)

	return &Player{
		id:   playerID,
		name: name,
	}
}

// func (player *Player) setPosition(x int, y int) {
// 	player.position.setPosition(x, y)
// }

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

func (player *Player) GetCurrentField() *Field {
	return player.currentField
}

func (player *Player) SetCurrentField(field *Field) {
	player.currentField = field
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

// func (player *Player) moveLeft() {
// 	player.currentField.horizontalFieldCode -= 1
// 	player.currentField.players = append(player.currentField.players, player)
// }

// func (player *Player) moveRight() {
// 	lastField := player.currentField

// 	player.currentField = NewField(lastField.verticalFieldCode, lastField.horizontalFieldCode+1)
// 	player.currentField.players = append(player.currentField.players, player)
// }

// func (player *Player) moveUp() {
// 	player.currentField.verticalFieldCode -= 1
// 	player.currentField.players = append(player.currentField.players, player)
// }

// func (player *Player) moveDown() {
// 	player.currentField.verticalFieldCode += 1
// 	player.currentField.players = append(player.currentField.players, player)
// }
