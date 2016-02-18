type GameMap struct {
  size int
}

func NewGameMap(size int) *GameMap {
    return &GameMap{
        size: size,
    }
}
