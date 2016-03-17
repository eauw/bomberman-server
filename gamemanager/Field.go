package gamemanager

import (
	"fmt"
	// "log"
	"strconv"
)

type Field struct {
	id       string
	row      int // horizontal
	column   int // vertical
	special  *Special
	wall     *Wall
	explodes bool
	players  map[string]*Player
	bombs    []*Bomb
}

func NewField(row int, column int) *Field {
	fieldID := strconv.Itoa(row) + strconv.Itoa(column)

	return &Field{
		id:       fieldID,
		row:      row,
		column:   column,
		explodes: false,
		players:  make(map[string]*Player),
		bombs:    make([]*Bomb, 0),
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

func (field *Field) setSpecial(powerType string) {
	field.special = NewSpecial(powerType)
}

func (field *Field) setWall(destructible bool) {
	if destructible {
		field.wall = NewWall(destructible)
	}
}

func (field *Field) addNewBomb(player *Player) *Bomb {
	bomb := NewBomb(player, field)
	field.bombs = append(field.bombs, bomb)

	return bomb
}

func (field *Field) removeBomb(bomb *Bomb) {
	index := -1

	for i := range field.bombs {
		if field.bombs[i] == bomb {
			index = i
		}
	}

	if index > 0 {
		slice1 := field.bombs[:index]
		slice2 := field.bombs[index+1:]

		newArray := append(slice1, slice2...)

		field.bombs = newArray
	}

}

// mach aus 1 -> 01 usw bis 10, ab dann normal
func cleanFieldNumber(number int) string {
	var n string

	if number < 10 {
		n = "0" + strconv.Itoa(number)
	}

	return n
}
