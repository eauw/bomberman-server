package main

import (
	"strconv"
)

type Field struct {
	id                  string
	horizontalFieldCode int // hCode
	verticalFieldCode   int // vCode
	containsSpecial     bool
	containsWall        bool
}

func NewField(vCode int, hCode int) *Field {
	fieldID := strconv.Itoa(vCode) + strconv.Itoa(hCode)

	return &Field{
		id:                  fieldID,
		horizontalFieldCode: hCode,
		verticalFieldCode:   vCode,
		containsSpecial:     false,
		containsWall:        false,
	}
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
