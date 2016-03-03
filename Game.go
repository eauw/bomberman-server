package main

type Game struct {
  channel chan string
  gameMap GameMap
}

func NewGame() *Game {
  ch := make(chan string)
  return &Game{
    channel: ch,
    gameMap: NewGameMap(5)
  }
}
