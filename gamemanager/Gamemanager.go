package gamemanager

import (
	// "bomberman-server/tcpmessage"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"
)

type Manager struct {
	game               *Game
	channel            chan *GameChannelMessage
	mainChannel        chan string
	currentPlayerIndex int
	playersOrder       []string
	playersConn        map[string]net.Conn
	commandTimeout     int
}

func NewManager() *Manager {
	return &Manager{
		playersOrder:       []string{}, // hält die IDs der Spieler in einer zufälligen Reihenfolge
		playersConn:        map[string]net.Conn{},
		currentPlayerIndex: 0,
		commandTimeout:     1,
	}
}

func (manager *Manager) SetMainChannel(ch chan string) {
	manager.mainChannel = ch
}

func (manager *Manager) Start(rounds int, xSize int, ySize int) {
	// manager.generatePlayersOrder()

	// manager.currentPlayerIndex = 0

	// currentPlayer := manager.GetCurrentPlayer()

	// manager.mainChannel <- fmt.Sprintf("first player: %s", currentPlayer.id)
	// manager.notifyCurrentPlayer()

	manager.game = NewGame(xSize, ySize)
	log.Println(manager.GameState())
}

func (manager *Manager) GameStart() {
	manager.game.start()
}

func (manager *Manager) GetCurrentPlayer() *Player {
	currentPlayer := manager.game.getPlayerByID(manager.playersOrder[manager.currentPlayerIndex])
	return currentPlayer
}

func (manager *Manager) setNextPlayer() {
	if manager.currentPlayerIndex == len(manager.game.players)-1 {
		manager.currentPlayerIndex = 0
	} else {
		manager.currentPlayerIndex += 1
	}

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
	newPlayer.addBomb()
	manager.game.addPlayer(newPlayer)

	manager.playersConn[newPlayer.id] = conn

	return newPlayer
}

func (manager *Manager) GameState() string {
	gameMap := manager.game.gameMap.toString()

	infos := "\n"
	infos += fmt.Sprintf("Runde: %d, ", manager.game.currentRound)
	infos += fmt.Sprintf("Spieleranzahl: %d, ", len(manager.game.players))
	infos += fmt.Sprintf("Spielfeldgröße: x %d y %d, ", manager.game.gameMap.xSize, manager.game.gameMap.ySize)
	infos += fmt.Sprintf("Cmd-Timeout: %d, ", manager.commandTimeout)
	infos += "\n"

	gameState := "\n"
	gameState += infos
	gameState += "\n"
	gameState += gameMap
	gameState += "\n"

	return gameState
}

func (manager *Manager) MessageReceived(message string, player *Player) {
	if manager.game.started {

		messageSlice := strings.Split(message, "")

		if len(messageSlice) > 0 {
			// prüfen ob Spieler eine Bombe werfen will
			if messageSlice[0] == "b" {
				// prüfen ob Spieler aktuell überhaupt verfügbare Bomben hat
				available := 0
				for _, b := range player.bombs {
					if b.isPlaced == false {
						available += 1
					}
				}
				if available > 0 {
					field := manager.destinationField(player, messageSlice)
					manager.game.PlayerPlacesBomb(player, field)
					manager.gameStateRequestedByPlayer(player)
				}

			}
		}

		switch message {
		case "d":
			manager.game.PlayerMovesToRight(player)
			manager.gameStateRequestedByPlayer(player)
			break

		case "a":
			manager.game.PlayerMovesToLeft(player)
			manager.gameStateRequestedByPlayer(player)
			break

		case "w":
			manager.game.PlayerMovesToUp(player)
			manager.gameStateRequestedByPlayer(player)
			break

		case "s":
			manager.game.PlayerMovesToDown(player)
			manager.gameStateRequestedByPlayer(player)
			break

		case "x":
			manager.game.ExplodePlayersBombs(player)
			manager.gameStateRequestedByPlayer(player)
			break

		case "l":
			manager.gameStateRequestedByPlayer(player)

		case "n":
			manager.setNextPlayer()
			break
		}
		// } else {
		// 	conn := manager.playersConn[player.id]
		// 	conn.Write([]byte("nyt: not your turn!\n"))
		// }
	} else {
		conn := manager.playersConn[player.id]
		conn.Write([]byte("Game waiting for more players.\n"))
	}
}

func (manager *Manager) gameStateRequestedByPlayer(p *Player) {
	conn := manager.playersConn[p.id]
	conn.Write([]byte(manager.GameState()))
}

func (manager *Manager) notifyCurrentPlayer() {
	currentPlayer := manager.GetCurrentPlayer()
	if currentPlayer != nil {
		conn := manager.playersConn[currentPlayer.id]
		conn.Write([]byte("yt: Your turn\n"))
	}

}

// Gibt für einen gegebenen Spieler und ein Ziel das entsprechende Feld zurück.
func (manager *Manager) destinationField(player *Player, destination []string) *Field {
	distance, _ := strconv.Atoi(destination[1])
	direction := destination[2]

	pRow := player.currentField.row
	pCol := player.currentField.column

	var destinationField *Field

	// Prüfen ob der Player weiter werfen will als er darf
	if distance > player.throwrange {
		distance = player.throwrange
	}

	switch direction {
	// Norden
	case "n":
		destinationField, _ = manager.game.gameMap.getField(pRow-distance, pCol)
		for destinationField == nil {
			distance -= 1
			destinationField, _ = manager.game.gameMap.getField(pRow-distance, pCol)
		}

		break

	// Osten
	case "o":
		destinationField, _ = manager.game.gameMap.getField(pRow, pCol+distance)
		for destinationField == nil {
			distance -= 1
			destinationField, _ = manager.game.gameMap.getField(pRow, pCol+distance)
		}
		break

	// Süden
	case "s":
		destinationField, _ = manager.game.gameMap.getField(pRow+distance, pCol)
		for destinationField == nil {
			distance -= 1
			destinationField, _ = manager.game.gameMap.getField(pRow+distance, pCol)
		}
		break

	// Westen
	case "w":
		destinationField, _ = manager.game.gameMap.getField(pRow, pCol-distance)
		for destinationField == nil {
			distance -= 1
			destinationField, _ = manager.game.gameMap.getField(pRow, pCol-distance)
		}
		break
	}

	return destinationField
}
