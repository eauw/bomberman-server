package main

type Player struct {
  name string
  points int
}

func NewPlayer(name string) *Player {
  return &Player{
      name: name,
  }
}
