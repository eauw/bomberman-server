package gamemanager

import (
	"fmt"
	// "log"
	"strconv"
)

type Field struct {
	id              string
	row             int // horizontal
	column          int // vertical
	containsSpecial bool
	containsWall    bool
	players         map[string]*Player
}

func NewField(row int, column int) *Field {
	fieldID := strconv.Itoa(row) + strconv.Itoa(column)

	return &Field{
		id:              fieldID,
		row:             row,
		column:          column,
		containsSpecial: false,
		containsWall:    false,
		players:         make(map[string]*Player),
	}
}

func (field *Field) toString() string {
	i := "nil"
	j := "nil"

	i = strconv.Itoa(field.row)

	j = strconv.Itoa(field.column)

	fieldString := fmt.Sprintf("i%s j%s", i, j)

	return fieldString

}

// func (field *Field) addPlayer(player *Player) {
// 	field.players = append(field.players, player)
// 	player.currentField = field
// }

func (field *Field) addPlayer(player *Player) {
	field.players[player.id] = player
	player.currentField = field
}

func (field *Field) removePlayer(player *Player) {
	delete(field.players, player.id)
}

func (field *Field) setSpecial(b bool) {
	field.containsSpecial = b
}

func (field *Field) setWall(b bool) {
	field.containsWall = b
}

// mach aus 1 -> 01 usw bis 10, ab dann normal
func cleanFieldNumber(number int) string {
	var n string

	if number < 10 {
		n = "0" + strconv.Itoa(number)
	}

	return n
}
