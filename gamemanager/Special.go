package gamemanager

import "github.com/eauw/bomberman-server/helper"
import "strings"

/*

powerTypes:
r: Bombenstrahlreichweite erhöhen
b: Anzahl der Bomben erhöhen
h: Schutz vor Bombenstrahl

*/

type Special struct {
	powerType string
}

func NewSpecial(powerType string) *Special {
	specials := specials()

	if strings.Contains(specials, powerType) {
		return &Special{
			powerType: powerType,
		}
	} else {
		return nil
	}

}

func RandomSpecial() *Special {
	specials := strings.Split(specials(), ",")

	i := helper.RandomNumber(0, len(specials))

	return NewSpecial(specials[i])
}

func specials() string {
	return "r,b,h"
}
