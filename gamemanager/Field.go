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
	players  []*Player
	bombs    []*Bomb
}

func NewField(row int, column int) *Field {
	fieldID := strconv.Itoa(row) + strconv.Itoa(column)

	return &Field{
		id:       fieldID,
		row:      row,
		column:   column,
		explodes: false,
		players:  make([]*Player, 0), //make(map[string]*Player),
		bombs:    make([]*Bomb, 0),
	}
}

func (field *Field) toString() string {
	i := "nil"
	j := "nil"

	i = strconv.Itoa(field.row)

	j = strconv.Itoa(field.column)

	fieldString := fmt.Sprintf("x:%sy:%s", j, i)

	return fieldString

}

func (field *Field) addPlayer(player *Player) {
	field.players = append(field.players, player)
	player.currentField = field
}

func (field *Field) removePlayer(player *Player) {
	index := -1

	for i := range field.players {
		if field.players[i].id == player.id {
			index = i
		}
	}

	if index > -1 {
		slice1 := field.players[:index]
		slice2 := field.players[index+1:]

		newArray := append(slice1, slice2...)

		field.players = newArray
	}
}

func (field *Field) setSpecial(powerType string) {
	field.special = NewSpecial(powerType)
}

func (field *Field) setWall(destructible bool) {
	if destructible {
		field.wall = NewWall(destructible)
	}
}

func (field *Field) addBomb(bomb *Bomb) {
	field.bombs = append(field.bombs, bomb)
}

func (field *Field) removeBomb(bomb *Bomb) {
	field.bombs = RemoveBomb(field.bombs, bomb)
}

// mach aus 1 -> 01 usw bis 10, ab dann normal
func cleanFieldNumber(number int) string {
	var n string

	if number < 10 {
		n = "0" + strconv.Itoa(number)
	}

	return n
}
