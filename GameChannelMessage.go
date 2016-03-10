package main

import (
	"fmt"
)

type GameChannelMessage struct {
	text   string
	player *Player
	game   *Game
}

func NewGameChannelMessage(text string, player *Player, game *Game) *GameChannelMessage {

	return &GameChannelMessage{
		text:   text,
		player: player,
		game:   game,
	}
}

func NewGameChannelMessageFromTCPMessage(tcpMessage *TCPMessage, game *Game) *GameChannelMessage {
	p := game.getPlayerByIP(tcpMessage.senderIP)

	return &GameChannelMessage{
		text:   tcpMessage.text,
		player: p,
		game:   game,
	}
}

func (gcm *GameChannelMessage) toString() {
	mainChannel <- fmt.Sprintf("text: %s, senderIP: %s", gcm.text, gcm.player.ip)
}
