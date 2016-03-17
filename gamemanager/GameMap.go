package gamemanager

import "bomberman-server/helper"

// import "fmt"

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

	// arr := []*Field{newField1, newField2, newField3, newField4, newField6, newField7, newField8, newField9}

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

	// arr := []*Field{newField2, newField4, newField6, newField8}

	return arr
}

func (gameMap *GameMap) addBomb(b *Bomb) {
	gameMap.bombs = append(gameMap.bombs, b)
}

func (gameMap *GameMap) removeBomb(b *Bomb) {
	index := -1

	for i := range gameMap.bombs {
		if gameMap.bombs[i] == b {
			index = i
		}
	}

	if index > 0 {
		slice1 := gameMap.bombs[:index]
		slice2 := gameMap.bombs[index+1:]

		newArray := append(slice1, slice2...)

		gameMap.bombs = newArray
	}
}

func createFields(xSize int, ySize int) [][]*Field {
	//horizontalFieldCodes := []string{"a","b","c","d","e","f","g","h","i","j","k","l","m","n","o","p","q","r","s","t","u","v","w","x","y","z"}
	//verticalFieldCodes := []string{"01","02","03","04","05","06","07","08","09","10","11","12","13","14","15","16","17","18","19","20","21","22","23","24","25","26"}

	fields := make([][]*Field, xSize)
	for i := range fields {
		fields[i] = make([]*Field, ySize)
		for j := range fields[i] {
			fields[i][j] = NewField(i, j)
		}
	}

	// place walls and specials on the game map
	walls := 20
	destructibleWalls := 5
	specials := 5

	// place walls
	for i := 0; i <= walls; i++ {
		randomRow := helper.RandomNumber(0, xSize)
		randomColumn := helper.RandomNumber(0, ySize)

		// TODO: prüfen ob auf dem Feld schon so ein Element liegt

		fields[randomRow][randomColumn].wall = NewWall(true)
	}

	// place destructible walls
	for i := 0; i <= destructibleWalls; i++ {
		randomRow := helper.RandomNumber(0, xSize)
		randomColumn := helper.RandomNumber(0, ySize)

		// TODO: prüfen ob auf dem Feld schon so ein Element liegt

		fields[randomRow][randomColumn].wall = NewWall(false)
	}

	// place specials
	for i := 0; i <= specials; i++ {
		randomRow := helper.RandomNumber(0, xSize)
		randomColumn := helper.RandomNumber(0, ySize)

		// TODO: prüfen ob auf dem Feld schon so ein Element liegt

		fields[randomRow][randomColumn].special = RandomSpecial()
	}

	return fields
}

func (gm *GameMap) toString() string {
	mapString := "\n"
	// fmt.Println()
	for i := range gm.fields {
		for j := range gm.fields[i] {
			f := gm.fields[i][j]
			// h := f.horizontalFieldCode
			// v := f.verticalFieldCode

			if len(f.players) > 0 {
				// print players
				if len(f.players) > 1 {
					mapString += "P"
				} else {
					mapString += "p"
				}

			} else if len(f.bombs) > 0 {
				// print bombs

				mapString += "B"

			} else if f.wall != nil {
				// print walls
				if f.wall.isDestructible {
					mapString += "w"
				} else {
					mapString += "W"
				}

			} else if f.special != nil {
				// print specials
				mapString += f.special.powerType

			} else {
				//fmt.Printf("_") //fmt.Printf("i %d, j %d", h, v) //fmt.Print(h + v)
				mapString += "_"
			}

			// fmt.Print("|")
			mapString += "|"
		}
		// fmt.Println()
		mapString += "\n"
	}

	//gm.game.mainChannel <- mapString

	// for i := range gm.fields {
	// 	for j := range gm.fields[i] {
	// 		// f := gm.fields[i][j]
	// 		mapString += fmt.Sprintf("%d %d|", i, j)
	// 	}
	// 	mapString += "\n"
	// }

	return mapString
}
