package gamemanager

import (
	// "bomberman-server/tcpmessage"
	"fmt"
	// "log"
	"math/rand"
	"net"
	"time"
)

type Manager struct {
	game          *Game
	channel       chan *GameChannelMessage
	mainChannel   chan string
	currentPlayer *Player
	playersOrder  []string
	playersConn   map[string]net.Conn
}

func NewManager() *Manager {
	return &Manager{
		game:         NewGame(),
		playersOrder: []string{},
		playersConn:  map[string]net.Conn{},
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
	manager.notifyCurrentPlayer()
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

func (manager *Manager) PlayerConnected(ip string, conn net.Conn) *Player {
	newPlayer := NewPlayer("New Player", manager.game.gameMap.fields[0][0])
	newPlayer.SetIP(ip)
	manager.game.addPlayer(newPlayer)

	manager.playersConn[newPlayer.id] = conn

	return newPlayer
}

func (manager *Manager) GameState() string {
	gameState := manager.game.gameMap.toString()
	gameState += "\n"
	gameState += fmt.Sprintf("Runde: %d\n", manager.game.currentRound)

	return gameState
}

func (manager *Manager) MessageReceived(message string, player *Player) {
	if player.id == manager.currentPlayer.id {
		switch message {
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

		case "l":
			manager.gameStateRequestedByPlayer(player)

		case "end":
			break
		}
	} else {
		conn := manager.playersConn[player.id]
		conn.Write([]byte("not your turn!"))
	}
}

func (manager *Manager) gameStateRequestedByPlayer(p *Player) {
	conn := manager.playersConn[p.id]
	conn.Write([]byte(manager.GameState()))
}

func (manager *Manager) notifyCurrentPlayer() {
	if manager.currentPlayer != nil {
		conn := manager.playersConn[manager.currentPlayer.id]
		conn.Write([]byte("Your turn\n"))
	}

}
