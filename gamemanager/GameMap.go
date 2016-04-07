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

// Liefert ein Feld anhand der Indizes zurück
func (gameMap *GameMap) getField(row int, column int) *Field {
	if row < 0 || row > len(gameMap.fields) || column < 0 || column > len(gameMap.fields) {
		return nil
	} else {
		return gameMap.fields[row][column]
	}

}

// Liefert alle angrenzenden Felder eines Feldes als Array zurück. Das Ausgangsfeld ist nicht inbegriffen.
func (gameMap *GameMap) GetFieldsAroundField(f *Field) []*Field {
	arr := []*Field{}

	newField1 := gameMap.getField(f.row-1, f.column-1)
	if newField1 != nil {
		arr = append(arr, newField1)
	}

	newField2 := gameMap.upperFieldOfField(f)
	if newField2 != nil {
		arr = append(arr, newField2)
	}

	newField3 := gameMap.getField(f.row-1, f.column+1)
	if newField3 != nil {
		arr = append(arr, newField3)
	}

	newField4 := gameMap.leftFieldOfField(f)
	if newField4 != nil {
		arr = append(arr, newField4)
	}

	newField6 := gameMap.rightFieldOfField(f)
	if newField6 != nil {
		arr = append(arr, newField6)
	}

	newField7 := gameMap.getField(f.row+1, f.column-1)
	if newField7 != nil {
		arr = append(arr, newField7)
	}

	newField8 := gameMap.lowerFieldOfField(f)
	if newField8 != nil {
		arr = append(arr, newField8)
	}

	newField9 := gameMap.getField(f.row+1, f.column+1)
	if newField9 != nil {
		arr = append(arr, newField9)
	}

	return arr
}

// liefert das obere, rechte, untere und linke Feld eines Feldes als Array zurück. Das Ausgangsfeld ist nicht inbegriffen.
func (gameMap *GameMap) GetNOSWFieldsOfField(f *Field) []*Field {
	arr := []*Field{}

	newField2 := gameMap.upperFieldOfField(f)
	if newField2 != nil {
		arr = append(arr, newField2)
	}

	newField4 := gameMap.leftFieldOfField(f)
	if newField4 != nil {
		arr = append(arr, newField4)
	}

	newField6 := gameMap.rightFieldOfField(f)
	if newField6 != nil {
		arr = append(arr, newField6)
	}

	newField8 := gameMap.lowerFieldOfField(f)
	if newField8 != nil {
		arr = append(arr, newField8)
	}

	return arr
}

func (gameMap *GameMap) GetNOSWFieldsOfFieldWithReach(f *Field, reach int) []*Field {
	arr := []*Field{}

	field := f

	for i := 1; i <= reach; i++ {

		field = gameMap.upperFieldOfField(field)
		if field != nil {
			arr = append(arr, field)
		}
	}

	field = f

	for i := 1; i <= reach; i++ {

		field = gameMap.leftFieldOfField(field)
		if field != nil {
			arr = append(arr, field)
		}
	}

	field = f

	for i := 1; i <= reach; i++ {

		field = gameMap.rightFieldOfField(field)
		if field != nil {
			arr = append(arr, field)
		}
	}

	field = f

	for i := 1; i <= reach; i++ {

		field = gameMap.lowerFieldOfField(field)
		if field != nil {
			arr = append(arr, field)
		}
	}

	return arr
}

func (gameMap *GameMap) leftFieldOfField(f *Field) *Field {
	return gameMap.getField(f.row, f.column-1)
}

func (gameMap *GameMap) upperFieldOfField(f *Field) *Field {
	return gameMap.getField(f.row-1, f.column)
}

func (gameMap *GameMap) rightFieldOfField(f *Field) *Field {
	return gameMap.getField(f.row, f.column+1)
}

func (gameMap *GameMap) lowerFieldOfField(f *Field) *Field {
	return gameMap.getField(f.row+1, f.column)
}

func (gameMap *GameMap) addBomb(b *Bomb) {
	gameMap.bombs = append(gameMap.bombs, b)
}

func (gameMap *GameMap) removeBomb(b *Bomb) {
	gameMap.bombs = RemoveBomb(gameMap.bombs, b)

}

func createFields(xSize int, ySize int) [][]*Field {
	fieldsCount := xSize * ySize

	fields := make([][]*Field, xSize)
	for i := range fields {
		fields[i] = make([]*Field, ySize)

		for j := range fields[i] {
			field := NewField(i, j)
			fields[i][j] = field

			// Rand um das Spielfeld anlegen
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
	wallsCount := int(float64(fieldsCount) * 0.25)
	destructibleWallsCount := int(float64(fieldsCount) * 0.05)
	specialsCount := int(float64(fieldsCount) * 0.1)

	// place walls
	for i := 0; i <= wallsCount; i++ {
		randomRow := helper.RandomNumber(0, xSize)
		randomColumn := helper.RandomNumber(0, ySize)

		if fields[randomRow][randomColumn].wall == nil {
			fields[randomRow][randomColumn].wall = NewWall(true)
		}
	}

	// place destructible walls
	for i := 0; i <= destructibleWallsCount; i++ {
		randomRow := helper.RandomNumber(0, xSize)
		randomColumn := helper.RandomNumber(0, ySize)

		if fields[randomRow][randomColumn].wall == nil {
			fields[randomRow][randomColumn].wall = NewWall(false)
		}
	}

	// place specials
	for i := 0; i <= specialsCount; i++ {
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
