package gamemanager

import "bomberman-server/helper"

import "fmt"
import "github.com/fatih/color"

type GameMap struct {
	game   *Game
	xSize  int
	ySize  int
	fields [][]*Field
	bombs  []*Bomb
}

func NewGameMap(xSize int, ySize int) *GameMap {
	f := createFields(xSize, ySize)

	return &GameMap{
		xSize:  xSize,
		ySize:  ySize,
		fields: f,
		bombs:  []*Bomb{},
	}
}

// Liefert ein Feld anhand der Indizes zurück und True bei Erfolg, sonst False.
func (gameMap *GameMap) getField(row int, column int) (*Field, bool) {
	if row < 0 || row > len(gameMap.fields) || column < 0 || column > len(gameMap.fields) {
		return nil, false
	} else {
		return gameMap.fields[row][column], true
	}

}

// Liefert alle angrenzenden Felder eines Feldes als Array zurück. Das Ausgangsfeld ist nicht inbegriffen.
func (gameMap *GameMap) GetFieldsAroundField(f *Field) []*Field {
	arr := []*Field{}

	if newField1, x := gameMap.getField(f.row-1, f.column-1); x {
		arr = append(arr, newField1)
	}
	if newField2, x := gameMap.getField(f.row-1, f.column); x {
		arr = append(arr, newField2)
	}
	if newField3, x := gameMap.getField(f.row-1, f.column+1); x {
		arr = append(arr, newField3)
	}
	if newField4, x := gameMap.getField(f.row, f.column-1); x {
		arr = append(arr, newField4)
	}
	if newField6, x := gameMap.getField(f.row, f.column+1); x {
		arr = append(arr, newField6)
	}
	if newField7, x := gameMap.getField(f.row+1, f.column-1); x {
		arr = append(arr, newField7)
	}
	if newField8, x := gameMap.getField(f.row+1, f.column); x {
		arr = append(arr, newField8)
	}
	if newField9, x := gameMap.getField(f.row+1, f.column+1); x {
		arr = append(arr, newField9)
	}

	return arr
}

// liefert das obere, rechte, untere und linke Feld eines Feldes als Array zurück. Das Ausgangsfeld ist nicht inbegriffen.
func (gameMap *GameMap) GetNOSWFieldsOfField(f *Field) []*Field {
	arr := []*Field{}

	if newField2, x := gameMap.getField(f.row-1, f.column); x {
		arr = append(arr, newField2)
	}
	if newField4, x := gameMap.getField(f.row, f.column-1); x {
		arr = append(arr, newField4)
	}
	if newField6, x := gameMap.getField(f.row, f.column+1); x {
		arr = append(arr, newField6)
	}
	if newField8, x := gameMap.getField(f.row+1, f.column); x {
		arr = append(arr, newField8)
	}

	return arr
}

func (gameMap *GameMap) addBomb(b *Bomb) {
	gameMap.bombs = append(gameMap.bombs, b)
}

func (gameMap *GameMap) removeBomb(b *Bomb) {
	index := -1

	if len(gameMap.bombs) > 1 {
		for i, bomb := range gameMap.bombs {
			if bomb.id == b.id {
				index = i
			}
		}

		if index > 0 {
			slice1 := gameMap.bombs[:index]
			slice2 := gameMap.bombs[index+1:]

			newArray := append(slice1, slice2...)

			gameMap.bombs = newArray

		}
	} else {
		gameMap.bombs = []*Bomb{}
	}

}

func createFields(xSize int, ySize int) [][]*Field {
	fields := make([][]*Field, xSize)
	for i := range fields {
		fields[i] = make([]*Field, ySize)

		for j := range fields[i] {
			field := NewField(i, j)
			fields[i][j] = field
			if i == 0 {
				field.wall = NewWall(false)
			}
			if i == xSize-1 {
				field.wall = NewWall(false)
			}
			if j == 0 {
				field.wall = NewWall(false)
			}
			if j == ySize-1 {
				field.wall = NewWall(false)
			}
		}
	}

	// place walls and specials on the game map
	walls := 20
	destructibleWalls := 5
	specials := 10

	// place walls
	for i := 0; i <= walls; i++ {
		randomRow := helper.RandomNumber(0, xSize)
		randomColumn := helper.RandomNumber(0, ySize)

		if fields[randomRow][randomColumn].wall == nil {
			fields[randomRow][randomColumn].wall = NewWall(true)
		}
	}

	// place destructible walls
	for i := 0; i <= destructibleWalls; i++ {
		randomRow := helper.RandomNumber(0, xSize)
		randomColumn := helper.RandomNumber(0, ySize)

		if fields[randomRow][randomColumn].wall == nil {
			fields[randomRow][randomColumn].wall = NewWall(false)
		}
	}

	// place specials
	for i := 0; i <= specials; i++ {
		randomRow := helper.RandomNumber(0, xSize)
		randomColumn := helper.RandomNumber(0, ySize)

		field := fields[randomRow][randomColumn]
		if field.special == nil {
			if field.wall == nil {
				field.special = RandomSpecial()
			} else {
				if field.wall.isDestructible {
					field.special = RandomSpecial()
				}
			}
		}
	}

	return fields
}

func (gm *GameMap) toString() string {
	mapString := "\n"

	for i := range gm.fields {
		for j := range gm.fields[i] {
			f := gm.fields[i][j]

			fieldChar := "_"

			if f.explodes {

				fieldChar = "*"

			} else if len(f.players) > 0 {
				// print players
				if len(f.players) > 1 {
					fieldChar = "P"
				} else if f.players[0].isFox > 0 {
					// Fuchs
					fieldChar = "f"
				} else {
					// normaler Spieler
					fieldChar = "p"
				}

			} else if len(f.bombs) > 0 {
				// print bombs

				fieldChar = "B"

			} else if f.wall != nil {
				// print walls
				if f.wall.isDestructible {
					fieldChar = "w"
				} else {
					fieldChar = "W"
				}

			} else if f.special != nil {
				// print specials
				fieldChar = f.special.powerType

			}

			mapString += fieldChar + "|"
		}

		mapString += "\n"
	}

	return mapString
}

func (gm *GameMap) toStringForServer() string {
	mapString := "\n"

	red := color.New(color.BgRed).SprintFunc()

	for i := range gm.fields {
		for j := range gm.fields[i] {
			f := gm.fields[i][j]

			fieldChar := "_"

			if len(f.players) > 0 {
				// print players
				if len(f.players) > 1 {
					fieldChar = "P"
				} else if f.players[0].isFox > 0 {
					// Fuchs
					fieldChar = "f"
				} else {
					// normaler Spieler
					fieldChar = "p"
				}

			} else if len(f.bombs) > 0 {
				// print bombs

				fieldChar = "B"

			} else if f.wall != nil {
				// print walls
				if f.wall.isDestructible {
					fieldChar = "w"
				} else {
					fieldChar = "W"
				}

			} else if f.special != nil {
				// print specials
				fieldChar = f.special.powerType

			}
			//decorate with 'explodes'-state
			if f.explodes {

				fieldChar = fmt.Sprintf("%s", red(fieldChar))

			}

			mapString += fieldChar + "|"
		}

		mapString += "\n"
	}

	// mapString += "bombs:\n"

	// for _, b := range gm.bombs {
	// 	mapString += fmt.Sprintf("%s\n\n", b)
	// }

	// mapString += "\n"

	return mapString
}
