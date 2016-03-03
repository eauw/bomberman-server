type GameMap struct {
  size int
  fields [][]string
}

func NewGameMap(size int) *GameMap {
  var f = [size][size]string
  return &GameMap{
      size: size,
      fields: f,
  }
}
