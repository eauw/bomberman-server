package gamemanager

import (
	// "bomberman-server/tcpmessage"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Manager struct {
	game               *Game
	channel            chan *GameChannelMessage
	mainChannel        chan string
	specChannel        chan string
	currentPlayerIndex int
	playersOrder       []string
	playersConn        map[string]net.Conn
	commandTimeout     int
	rounds             []*Round
	mutex              *sync.Mutex
	currentRound       *Round
}

func NewManager() *Manager {
	return &Manager{
		playersOrder:       []string{}, // hält die IDs der Spieler in einer zufälligen Reihenfolge
		playersConn:        map[string]net.Conn{},
		currentPlayerIndex: 0,
		commandTimeout:     1,
		rounds:             []*Round{},
		mutex:              &sync.Mutex{},
	}
}

func (manager *Manager) SetMainChannel(ch chan string) {
	manager.mainChannel = ch
}

func (manager *Manager) SetSpecChannel(ch chan string) {
	manager.specChannel = ch
}

func (manager *Manager) Start(rounds int, xSize int, ySize int) {

	if rounds < 1 {
		rounds = 1
	}

	for i := 1; i <= rounds; i++ {
		round := NewRound()
		round.id = i
		manager.rounds = append(manager.rounds, round)
	}

	manager.game = NewGame(xSize, ySize)
	manager.currentRound = manager.rounds[0]
	log.Println(manager.GameState(manager.game.gameMap.toStringForServer()))
}

func (manager *Manager) GameStart() {
	manager.game.start()

	for _, p := range manager.game.players {
		manager.sendGameStateToPlayer(p)
	}

	log.Println(manager.GameState(manager.game.gameMap.toStringForServer()))

	manager.broadcastWaiting()
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
	newPlayer.name = "Player" + strconv.Itoa(len(manager.game.players))

	manager.playersConn[newPlayer.id] = conn

	return newPlayer
}

func (manager *Manager) GameState(mapString string) string {

	infos := "\n"
	if manager.game.started {
		infos += fmt.Sprintf("round:%d/%d,", manager.currentRound.id, len(manager.rounds))
	} else {
		infos += fmt.Sprintf("rounds:%d,", len(manager.rounds))
	}

	infos += fmt.Sprintf("players:%d,", len(manager.game.players))

	x := manager.game.gameMap.xSize
	y := manager.game.gameMap.ySize

	xString := strconv.Itoa(x)

	if x < 10 {
		xString = "0" + xString
	}

	yString := strconv.Itoa(y)

	if y < 10 {
		yString = "0" + yString
	}

	infos += fmt.Sprintf("mapsize:x%sy%s,", xString, yString)
	infos += fmt.Sprintf("timeout:%d,", manager.commandTimeout)
	infos += "\n"

	gameState := "\n"
	gameState += "***********************************************************"
	gameState += "\n"
	gameState += infos
	gameState += "\n"
	gameState += "map:"
	gameState += mapString
	gameState += "\n"

	// gamestatetable
	if manager.game.started {
		// Tabelle erstellen mit dem Fuchs an erster Stelle
		playersTable := make([]*Player, len(manager.game.players))
		i := 1
		for _, p := range manager.game.players {
			if p.isFox > 0 {
				playersTable[0] = p
			} else {
				playersTable[i] = p
				i += 1
			}
		}

		gameStateTable := "scoretable:\n"
		for _, p := range playersTable {
			gameStateTable += fmt.Sprintf("name:%s,score:%d,field:%s;\n", p.name, p.points, p.currentField.toString())
		}
		gameStateTable += "/scoretable"
		gameState += gameStateTable
	}
	gameState += "\n"

	// gameState += "bombs:\n"

	// for _, p := range manager.game.players {
	// 	gameState += fmt.Sprintf("player: %s, bombs: %s\n", p.name, p.bombs)
	// }

	gameState += "***********************************************************"
	gameState += "\n"

	return gameState
}

func (manager *Manager) MessageReceived(message string, player *Player) {
	// mutex.Lock()
	// defer mutex.Unlock()

	log.Printf("Message >%s< received from player >%s<", message, player.name)
	conn := manager.playersConn[player.id]

	// mit q verlässt der Spieler den Server
	if message == "q" {
		manager.playerQuit(player)
	} else {

		// Befehle werden erst entgegengenommen wenn das Spiel gestartet wurde
		if manager.game.started {

			playerCommands := manager.currentRound.playerCommands

			if _, alreadyExits := playerCommands[player.id]; alreadyExits {
				conn.Write([]byte("Your already have send a message.\n"))
			} else {
				manager.currentRound.playerCommands[player.id] = message
			}

			if len(playerCommands) == len(manager.game.players) {
				manager.ProcessRound(manager.currentRound)
			}

		} else {
			if strings.Contains(message, "name:") {
				name := strings.TrimPrefix(message, "name:")
				player.SetName(name)
				conn.Write([]byte("YourID:"))
				conn.Write([]byte(player.GetID()))
				conn.Write([]byte("\n"))
				conn.Write([]byte("YourName:"))
				conn.Write([]byte(player.GetName()))
				conn.Write([]byte("\n"))
			} else {
				conn.Write([]byte("Game waiting for more players.\n"))
			}
		}
	}
}

func (manager *Manager) ProcessRound(round *Round) {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	fields := manager.game.gameMap.fields
	bombs := manager.game.gameMap.bombs

	log.Printf("Processing Round %d\n", round.id)

	for playerID, command := range round.playerCommands {
		player := manager.game.getPlayerByID(playerID)

		messageSlice := strings.Split(command, "")

		// prüfen ob Spieler eine Bombe werfen will
		if len(messageSlice) > 0 {
			if messageSlice[0] == "b" {
				// prüfen ob Spieler aktuell überhaupt verfügbare Bomben hat
				available := 0
				for _, b := range player.bombs {
					if b.field == nil {
						available += 1
					}
				}
				if available > 0 {
					// Bomben sind verfügbar
					field := manager.destinationField(player, messageSlice)
					manager.game.PlayerPlacesBomb(player, field)
				}
			}
		}

		switch command {
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

		// case "x":
		// 	manager.game.ExplodePlayersBombs(player)
		// 	break

		case "l":
			manager.sendGameStateToPlayer(player)
			break

		case "n":
			// nothing
			manager.sendGameStateToPlayer(player)
			break

		}
	}

	// Bomben Timer runterzählen und ggf. explodieren lassen
	for _, b := range bombs {
		b.timer -= 1
		if b.field != nil && b.timer <= 0 {
			b.explode(manager.game.gameMap)
		}
	}

	// Punkte des Fuchses erhöhen und Schutz abziehen falls nötig
	for _, p := range manager.game.players {
		if p.isFox > 0 {
			p.isFox += 1
			p.points += p.isFox
		}
		if p.protection > 0 {
			p.protection -= 1
		}
		if p.isParalyzed > 0 {
			p.isParalyzed -= 1
		}
	}

	manager.broadcastGamestate()

	// nächste Runde setzen
	roundIdx := round.id

	// prüfen ob die letzte Runde erreicht ist
	if roundIdx >= len(manager.rounds) {
		// dann beenden
		manager.finishGame()
	} else {
		// sonst nächste Runde
		manager.currentRound = manager.rounds[roundIdx]
	}

	// Neue Runde konfigurieren

	// explodierte Bomben zurücksetzen
	for i := range fields {
		for j := range fields[i] {
			f := fields[i][j]
			f.explodes = false
		}
	}

	// TODO: Specials respawn implementieren

	manager.broadcastWaiting()

}

func (manager *Manager) broadcastWaiting() {
	log.Println("waiting for commands")
	for _, p := range manager.game.players {
		conn := manager.playersConn[p.id]
		conn.Write([]byte("wfyc: waiting for your command\n"))
	}
}

func (manager *Manager) broadcastGamestate() {
	for _, p := range manager.game.players {
		manager.sendGameStateToPlayer(p)
	}

	log.Println(manager.GameState(manager.game.gameMap.toStringForServer()))

	manager.specChannel <- manager.GameState(manager.game.gameMap.toStringForServer())
}

func (manager *Manager) sendGameStateToPlayer(p *Player) {
	conn := manager.playersConn[p.id]
	conn.Write([]byte(buildHeader(manager.GameState(manager.game.gameMap.toString())))) // ???
	conn.Write([]byte(manager.GameState(manager.game.gameMap.toString())))
}

// Gibt für einen gegebenen Spieler und ein Ziel das entsprechende Feld zurück.
func (manager *Manager) destinationField(player *Player, destination []string) *Field {

	var distance int
	var direction string

	if len(destination) < 3 {
		return player.currentField
	} else {
		distance, _ = strconv.Atoi(destination[1])
		direction = destination[2]
	}

	// prüfen ob Richtung gültig ist
	validDirections := "wasd"
	if strings.Contains(validDirections, direction) == false {
		return player.currentField
	}

	pRow := player.currentField.row
	pCol := player.currentField.column

	var destinationField *Field

	// Prüfen ob der Player weiter werfen will als er darf
	if distance > player.throwrange {
		distance = player.throwrange
	}

	switch direction {
	// Norden
	case "w":
		destinationField, _ = manager.game.gameMap.getField(pRow-distance, pCol)
		for destinationField == nil {
			distance -= 1
			destinationField, _ = manager.game.gameMap.getField(pRow-distance, pCol)
		}

		break

	// Osten
	case "d":
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
	case "a":
		destinationField, _ = manager.game.gameMap.getField(pRow, pCol-distance)
		for destinationField == nil {
			distance -= 1
			destinationField, _ = manager.game.gameMap.getField(pRow, pCol-distance)
		}
		break
	}

	return destinationField
}

func (manager *Manager) playerQuit(player *Player) {

	conn := manager.playersConn[player.id]
	conn.Write([]byte("good bye\n"))

	delete(manager.playersConn, player.id)
	manager.game.removePlayer(player)

	log.Printf("Player %s has left the game.\n", player.id)

	conn.Close()
}

func (manager *Manager) finishGame() {
	log.Println("Game is over.")

	for _, conn := range manager.playersConn {
		conn.Write([]byte("game is over\n"))
		conn.Close()
	}
}

func (manager *Manager) sendMessageToPlayer(message string, player *Player) {

}

func buildHeader(message string) string {
	messageBytes := []byte(message)
	messageLength := len(messageBytes)

	// header := &TCPHeader{bytes, messageLength}

	headerString := fmt.Sprintf("messageLength:%d", messageLength)

	return headerString
}

type TCPHeader struct {
	bytes         []byte
	messageLength int
}
