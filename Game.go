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
	s := "Players\n"
	for i, v := range game.players {
		playerCount := i+1
		s += fmt.Sprintf("%d. ID: %s | IP: %s | Name: %s \n", playerCount, v.id, v.ip, v.name)
	}

	return s
}

// func (game *Game) placePlayers() {
// 	for p := range game.players {
//
// 	}
// }
