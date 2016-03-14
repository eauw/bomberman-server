package gamemanager

import (
	"bomberman-server/tcpmessage"
)

type Manager struct {
	game        *Game
	channel     chan *GameChannelMessage
	mainChannel chan string
}

func NewManager() *Manager {
	return &Manager{
		game: NewGame(),
	}
}

func (manager *Manager) SetMainChannel(ch chan string) {
	manager.mainChannel = ch
}

func (manager *Manager) Start() {

}

func (manager *Manager) PlayerConnected(player *Player) {
	manager.game.addPlayer(player)
}

func (manager *Manager) GameState() string {
	return manager.game.gameMap.toString()
}

func (manager *Manager) MessageReceived(tcpMessage *tcpmessage.TCPMessage) {
	player := manager.game.getPlayerByIP(tcpMessage.SenderIP)
	switch tcpMessage.Text {
	case "move right":
		manager.game.PlayerMovesToRight(player)
		break

	case "move left":
		manager.game.PlayerMovesToLeft(player)
		break

	case "move up":
		manager.game.PlayerMovesToUp(player)
		break

	case "move down":
		manager.game.PlayerMovesToDown(player)
		break

	case "bomb":
		manager.game.PlayerPlacesBomb(player)
		break

	}
}
