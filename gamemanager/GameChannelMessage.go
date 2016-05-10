package gamemanager

import (
	// "fmt"
	"github.com/eauw/bomberman-server/tcpmessage"
	"time"
)

type GameChannelMessage struct {
	text      string
	player    *Player
	timeStamp time.Time
}

func NewGameChannelMessage(text string, player *Player) GameChannelMessage {

	return GameChannelMessage{
		text:      text,
		player:    player,
		timeStamp: time.Now(),
	}
}

func NewGameChannelMessageFromTCPMessage(tcpMessage *tcpmessage.TCPMessage, game *Game) *GameChannelMessage {
	p := game.getPlayerByIP(tcpMessage.SenderIP)

	return &GameChannelMessage{
		text:   tcpMessage.Text,
		player: p,
	}
}

func (self *GameChannelMessage) GetText() string {
	return self.text
}

func (self *GameChannelMessage) GetPlayer() *Player {
	return self.player
}

// func (gcm *GameChannelMessage) toString() {
// 	mainChannel <- fmt.Sprintf("text: %s, senderIP: %s", gcm.Text, gcm.Player.Ip)
// }
