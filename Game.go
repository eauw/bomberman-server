package main

import(
	"fmt"
)

type Game struct {
	channel chan string
	gameMap *GameMap
	players []*Player
}

func NewGame() *Game {
	ch := make(chan string)

	return &Game{
		channel: ch,
		gameMap: NewGameMap(20),
	}
}

func (game *Game) addPlayer(player *Player) {
	game.players = append(game.players, player)
}

func (game *Game) removePlayer(player *Player) {

}

func (game *Game) printPlayers() string {
	s := ""
	for _, v := range game.players {
		s += fmt.Sprintf("Player-ID: %s\n",v.id)
	}

	return s
}

// func (game *Game) placePlayers() {
// 	for p := range game.players {
//
// 	}
// }
