package gamemanager

import "bomberman-server/helper"

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

	counter := 0
	for v := range specials {
		if specials[v] == powerType {
			counter += 1
		}
	}

	if counter == 0 {
		return nil
	} else {
		return &Special{
			powerType: powerType,
		}
	}
}

func RandomSpecial() *Special {
	specials := specials()

	i := helper.RandomNumber(0, len(specials))

	return NewSpecial(specials[i])
}

func specials() []string {
	return []string{"r", "b", "h"}
}
