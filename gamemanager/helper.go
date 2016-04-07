package gamemanager

func RemoveBomb(bombs []*Bomb, b *Bomb) []*Bomb {
	index := -1

	if len(bombs) > 1 {
		for i, bomb := range bombs {
			if bomb.id == b.id {
				index = i
			}
		}

		if index >= 0 {
			slice1 := bombs[:index]
			slice2 := bombs[index+1:]

			newArray := append(slice1, slice2...)

			return newArray

		}
	} else {
		return []*Bomb{}
	}

	return bombs
}
