package gamemanager

import (
	// "fmt"
	"bomberman-server/tcpmessage"
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

func NewGameChannelMessageFromTCPMessage(tcpMessage *tcpmessage.TCPMessage, game *Game) *GameChannelMessage {
	p := game.getPlayerByIP(tcpMessage.SenderIP)

	return &GameChannelMessage{
		text:   tcpMessage.Text,
		player: p,
		game:   game,
	}
}

func (self *GameChannelMessage) GetText() string {
	return self.text
}

func (self *GameChannelMessage) GetPlayer() *Player {
	return self.player
}

func (self *GameChannelMessage) GetGame() *Game {
	return self.game
}

// func (gcm *GameChannelMessage) toString() {
// 	mainChannel <- fmt.Sprintf("text: %s, senderIP: %s", gcm.Text, gcm.Player.Ip)
// }
