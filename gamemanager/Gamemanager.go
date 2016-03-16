package gamemanager

import (
	"bomberman-server/tcpmessage"
	"fmt"
	// "log"
	"math/rand"
	"time"
)

type Manager struct {
	game          *Game
	channel       chan *GameChannelMessage
	mainChannel   chan string
	currentPlayer *Player
	playersOrder  []string
}

func NewManager() *Manager {
	return &Manager{
		game:         NewGame(),
		playersOrder: []string{},
	}
}

func (manager *Manager) SetMainChannel(ch chan string) {
	manager.mainChannel = ch
}

func (manager *Manager) Start() {
	manager.generatePlayersOrder()
	// log.Print(manager.playersOrder)
	manager.setCurrentPlayer(manager.game.getPlayerByID(manager.playersOrder[0]))

	manager.mainChannel <- fmt.Sprintf("first player: %s", manager.currentPlayer.id)
}

func (manager *Manager) GetCurrentPlayer() *Player {
	return manager.currentPlayer
}

func (manager *Manager) setCurrentPlayer(p *Player) {
	manager.currentPlayer = p
}

func (manager *Manager) PlayersCount() int {
	if manager.game.players != nil {
		return len(manager.game.players)
	} else {
		return 0
	}
}

// Erstellt eine zufällige Spielerreihenfolge für die Runden.
func (manager *Manager) generatePlayersOrder() {
	a := []string{}

	for i := range manager.game.players {
		a = append(a, manager.game.players[i].id)
	}

	rand.Seed(time.Now().UnixNano())

	// shuffle
	for i := range a {
		j := rand.Intn(i + 1)
		a[i], a[j] = a[j], a[i]
	}

	manager.playersOrder = a
}

func (manager *Manager) PlayerConnected(ip string) *Player {
	newPlayer := NewPlayer("New Player", manager.game.gameMap.fields[0][0])
	newPlayer.SetIP(ip)
	manager.game.addPlayer(newPlayer)

	return newPlayer
}

func (manager *Manager) GameState() string {
	gameState := manager.game.gameMap.toString()
	gameState += "\n"
	gameState += fmt.Sprintf("Runde: %d\n", manager.game.currentRound)
	return gameState
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
