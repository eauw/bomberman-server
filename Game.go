package main

import (
	"fmt"
	"log"
)

type Game struct {
	channel     chan *GameChannelMessage
	mainChannel chan string
	gameMap     *GameMap
	players     []*Player
}

func NewGame() *Game {
	ch := make(chan *GameChannelMessage)
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
		var gameChannelMessage = <-game.channel
		handleGameChannelMessage(gameChannelMessage)
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
		playerCount := i + 1
		s += fmt.Sprintf("%d. ID: %s | IP: %s | Name: %s \n", playerCount, v.id, v.ip, v.name)
	}

	return s
}

func (game *Game) getPlayerByIP(ip string) *Player {
	for _, p := range game.players {
		if p.ip == ip {
			return p
		}
	}

	return nil
}

func handleGameChannelMessage(gcm *GameChannelMessage) {
	log.Printf("aaaa: %s", gcm.player.toString())
	switch gcm.text {
	case "move up":
		gcm.player.moveUp()
		break

	case "move down":
		gcm.player.moveDown()
		break

	case "move left":
		gcm.player.moveLeft()
		break

	case "move right":
		gcm.player.moveRight()
		break

	}
}

// func (game *Game) placePlayers() {
// 	for p := range game.players {
//
// 	}
// }
