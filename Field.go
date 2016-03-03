package main

import (
	"strconv"
)

type Field struct {
	number string
	horizontalFieldCode int 	// hCode
	verticalFieldCode int	 		// vCode
}

func NewField(hCode int, vCode int) *Field {
	return &Field{
		horizontalFieldCode: hCode,
		verticalFieldCode: vCode,
	}
}

// mach aus 1 -> 01 usw bis 10, ab dann normal
func cleanFieldNumber(number int) string {
	var n string

	if number  < 10 {
		n = "0" + strconv.Itoa(number)
	}

	return n
}
