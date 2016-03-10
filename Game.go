package main

import(
	"fmt"
)

type Game struct {
	channel chan string
	mainChannel chan string
	gameMap *GameMap
	players []*Player
}

func NewGame() *Game {
	ch := make(chan string)
	gm := NewGameMap(20)

	newGame := &Game{
		channel: ch,
		gameMap: gm,
	}

	gm.game = newGame

	return newGame
}

func (game *Game) start() {
	go game.handleGameChannel()
}

// receives all information about the game
func (game *Game) handleGameChannel() {
	for {
		var x = <-game.channel
		//fmt.Printf("game channel: %s", x)
		switch x {
		case "":
			break

		}
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
