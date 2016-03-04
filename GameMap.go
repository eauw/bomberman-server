package main

import (
	"fmt"
)

type GameMap struct {
	size   int
	fields [][]*Field
}

func NewGameMap(size int) *GameMap {
	f := createFields(size)

	return &GameMap{
		size:   size,
		fields: f,
	}
}

func createFields(size int) [][]*Field {
	//horizontalFieldCodes := []string{"a","b","c","d","e","f","g","h","i","j","k","l","m","n","o","p","q","r","s","t","u","v","w","x","y","z"}
	//verticalFieldCodes := []string{"01","02","03","04","05","06","07","08","09","10","11","12","13","14","15","16","17","18","19","20","21","22","23","24","25","26"}

	f := make([][]*Field, size)
	for i := range f {
		f[i] = make([]*Field, size)
		for j := range f[i] {
			f[i][j] = NewField(i,j)
		}
	}

	return f
}

func (gm *GameMap) toString() {
	for i := range gm.fields {
		for j := range gm.fields[i] {
			h := gm.fields[i][j].horizontalFieldCode
			v := gm.fields[i][j].verticalFieldCode
			fmt.Print(h+v)
			fmt.Print("|")
		}
		fmt.Println()
	}
}
