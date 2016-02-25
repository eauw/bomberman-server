package main

type Game struct {
  channel chan string
}

func NewGame() *Game {
  ch := make(chan string)
  return &Game{
    channel: ch,
  }
}
