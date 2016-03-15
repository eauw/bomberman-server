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

func (manager *Manager) PlayerConnected(ip string) *Player {
	newPlayer := NewPlayer("New Player", manager.game.gameMap.fields[0][0])
	newPlayer.SetIP(ip)
	manager.game.addPlayer(newPlayer)

	return newPlayer
}

func (manager *Manager) GameState() string {
	return manager.game.gameMap.toString()
}

func (manager *Manager) MessageReceived(tcpMessage *tcpmessage.TCPMessage) {
	player := manager.game.getPlayerByIP(tcpMessage.SenderIP)
	switch tcpMessage.Text {
	case "d":
		manager.game.PlayerMovesToRight(player)
		break

	case "a":
		manager.game.PlayerMovesToLeft(player)
		break

	case "w":
		manager.game.PlayerMovesToUp(player)
		break

	case "s":
		manager.game.PlayerMovesToDown(player)
		break

	case "b":
		manager.game.PlayerPlacesBomb(player)
		break

	case "x":
		manager.game.ExplodeBomb()
		break

	case "end":
		break
	}
}
