package gamemanager

import (
	// "bomberman-server/tcpmessage"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

var timer *time.Timer

type Manager struct {
	currentGame    *Game
	games          []*Game
	channel        chan GameChannelMessage
	mainChannel    chan string
	specChannel    chan string
	playersConn    map[string]net.Conn
	commandTimeout float64
}

func NewManager() *Manager {
	ch := make(chan GameChannelMessage, 2)
	manager := &Manager{
		playersConn: map[string]net.Conn{},
		channel:     ch,
	}

	go manager.channelHandler()

	return manager
}

func (manager *Manager) GetGameChannel() chan GameChannelMessage {
	return manager.channel
}

func (manager *Manager) SetMainChannel(ch chan string) {
	manager.mainChannel = ch
}

func (manager *Manager) SetSpecChannel(ch chan string) {
	manager.specChannel = ch
}

func (manager *Manager) Start(rounds int, height int, width int, gamesCount int, timeout float64) {

	if rounds < 1 {
		rounds = 1
	}

	manager.commandTimeout = timeout

	// init games and rounds
	for i := 1; i <= gamesCount; i++ {
		newGame := NewGame(height, width)
		newGame.id = i

		for j := 1; j <= rounds; j++ {
			round := NewRound()
			round.id = j
			newGame.rounds = append(newGame.rounds, round)
		}

		manager.games = append(manager.games, newGame)
	}

	manager.currentGame = manager.games[0]
	manager.currentGame.currentRound = manager.currentGame.rounds[0]
	log.Println(manager.GameState(manager.currentGame.gameMap.toStringForServer()))
}

func (manager *Manager) GameStart() {
	manager.currentGame.start()

	for _, p := range manager.currentGame.players {
		manager.sendGameStateToPlayer(p)
	}

	log.Println(manager.GameState(manager.currentGame.gameMap.toStringForServer()))

	manager.broadcastWaiting()

	manager.timeout()
}

func (manager *Manager) timeout() {
	timer = time.NewTimer(time.Duration(float64(time.Second) * manager.commandTimeout))
	go func() {
		<-timer.C

		gameChannelMessage := NewGameChannelMessage("processRound", nil)
		manager.channel <- gameChannelMessage
	}()
}

func (manager *Manager) PlayersCount() int {
	if manager.currentGame.players != nil {
		return len(manager.currentGame.players)
	} else {
		return 0
	}
}

func (manager *Manager) PlayerConnected(ip string, conn net.Conn) *Player {
	newPlayer := NewPlayer("New Player")
	newPlayer.SetIP(ip)
	newPlayer.addBomb()
	manager.currentGame.addPlayer(newPlayer)
	newPlayer.name = "Player" + strconv.Itoa(len(manager.currentGame.players))

	manager.playersConn[newPlayer.id] = conn

	return newPlayer
}

func (manager *Manager) GameState(mapString string) string {

	infos := "\n"
	infos += fmt.Sprintf("game:%d/%d,", manager.currentGame.id, len(manager.games))
	if manager.currentGame.started {
		infos += fmt.Sprintf("round:%d/%d,", manager.currentGame.currentRound.id, len(manager.currentGame.rounds))
	} else {
		infos += fmt.Sprintf("rounds:%d,", len(manager.currentGame.rounds))
	}

	infos += fmt.Sprintf("players:%d,", len(manager.currentGame.players))

	x := manager.currentGame.gameMap.height
	y := manager.currentGame.gameMap.width

	xString := strconv.Itoa(x)

	if x < 10 {
		xString = "0" + xString
	}

	yString := strconv.Itoa(y)

	if y < 10 {
		yString = "0" + yString
	}

	infos += fmt.Sprintf("mapsize:x%sy%s,", xString, yString)
	infos += fmt.Sprintf("timeout:%.2fs,", manager.commandTimeout)
	infos += "\n"

	gameState := "\033[H\033[2J"
	gameState += "***********************************************************"
	gameState += "\n"
	gameState += infos
	gameState += "\n"
	gameState += "map:"
	gameState += mapString
	gameState += "\n"

	// gamestatetable
	if manager.currentGame.started {
		// Tabelle erstellen mit dem Fuchs an erster Stelle
		playersTable := make([]*Player, len(manager.currentGame.players))
		i := 1
		for _, p := range manager.currentGame.players {
			if p.isFox > 0 {
				playersTable[0] = p
			} else {
				playersTable[i] = p
				i += 1
			}
		}

		gameStateTable := "scoretable:\n"
		for _, p := range playersTable {
			gameStateTable += fmt.Sprintf("name:%s,score:%d,%s;\n", p.name, p.points, p.currentField.toString())
		}
		gameStateTable += "/scoretable"
		gameState += gameStateTable
	}
	gameState += "\n"

	// gameState += "bombs:\n"

	// for _, p := range manager.currentGame.players {
	// 	gameState += fmt.Sprintf("player: %s, bombs: %s\n", p.name, p.bombs)
	// }

	gameState += "***********************************************************"
	gameState += "\n"

	return gameState
}

func (manager *Manager) channelHandler() {
	for {
		gameChannelMessage := <-manager.channel

		if gameChannelMessage.text == "processRound" && gameChannelMessage.player == nil {
			manager.ProcessRound(manager.currentGame.currentRound)
		} else {
			manager.messageReceived(gameChannelMessage.text, gameChannelMessage.player, gameChannelMessage.timeStamp)
		}

	}
}

func (manager *Manager) messageReceived(message string, player *Player, timestamp time.Time) {
	log.Printf("Message >%s< received from player >%s<", message, player.name)
	conn := manager.playersConn[player.id]

	// mit q verlässt der Spieler den Server
	if message == "q" {
		manager.playerQuit(player)
	} else {

		// Befehle werden erst entgegengenommen wenn das Spiel gestartet wurde
		if manager.currentGame.started {

			playerCommands := manager.currentGame.currentRound.playerCommands

			if _, alreadyExits := playerCommands[player.id]; alreadyExits {
				conn.Write([]byte("You already have send a message.\n"))
			} else {
				manager.currentGame.currentRound.playerCommands[player.id] = PlayerCommand{message, timestamp}
			}

			if len(playerCommands) == len(manager.currentGame.players) {
				// manager.ProcessRound(manager.currentRound)
				if timer != nil {
					timer.Stop()

				}
				gameChannelMessage := NewGameChannelMessage("processRound", nil)
				manager.channel <- gameChannelMessage

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
	fields := manager.currentGame.gameMap.fields
	bombs := manager.currentGame.gameMap.bombs

	log.Printf("Processing Round %d\n", round.id)

	players := []*Player{}

	for playerID, command := range round.playerCommands {
		player := manager.currentGame.getPlayerByID(playerID)
		players = append(players, player)

		messageSlice := strings.Split(command.message, "")

		// prüfen ob Spieler eine Bombe werfen will
		if len(messageSlice) > 0 {
			if messageSlice[0] == "b" {
				// prüfen ob Spieler aktuell überhaupt verfügbare Bomben hat
				available := false
				for _, b := range player.bombs {
					if b.field == nil {
						// Bomben sind verfügbar
						field := manager.destinationField(player, messageSlice)
						manager.currentGame.PlayerPlacesBomb(player, field)
						available = true
						break;
					}
				}
				if !available {
					player.msg += "E: out of bombs\n"
				}
			}
		}

		switch command.message {
		case "d":
			manager.currentGame.PlayerMovesToRight(player)

		case "a":
			manager.currentGame.PlayerMovesToLeft(player)

		case "w":
			manager.currentGame.PlayerMovesToUp(player)

		case "s":
			manager.currentGame.PlayerMovesToDown(player)

		// case "x":
		// 	manager.currentGame.ExplodePlayersBombs(player)

		case "l":
			manager.sendGameStateToPlayer(player)

		case "n":
			// nothing
			manager.sendGameStateToPlayer(player)

		}
	}

	// Prüfen ob Fuchs gefangen wurde // NICHT GETESTET!
	var fox *Player
	for _, v := range players {

		if v.isFox > 0 {
			fox = v
		}

	}
	if len(fox.currentField.players) >= 2 {
		playersOnFoxFieldWithoutFox := []*Player{}
		for _, v := range fox.currentField.players {
			if v.isFox == 0 {
				playersOnFoxFieldWithoutFox = append(playersOnFoxFieldWithoutFox, v)
			}
		}

		// PlayerCommand mit frühestem timestamp finden
		var p = playersOnFoxFieldWithoutFox[0]               // ersten player als frühesten annehmen
		var timestamp = round.playerCommands[p.id].timestamp // timestamp von erstem player festhalten
		for _, v := range playersOnFoxFieldWithoutFox {
			playerCommand := round.playerCommands[v.id]
			if playerCommand.timestamp.Sub(timestamp) < 0 {
				//playerCommand ist früher
				timestamp = playerCommand.timestamp
				p = v
			}
		}

		fox.isFox = 0
		p.isFox += 1
		manager.currentGame.teleportPlayer(p)

	}

	// Bomben Timer runterzählen und ggf. explodieren lassen
	for _, b := range bombs {
		b.timer -= 1
		if b.field != nil && b.timer <= 0 {
			b.explode(manager.currentGame.gameMap)
		}
	}

	// Punkte des Fuchses erhöhen und Schutz abziehen falls nötig
	for _, p := range manager.currentGame.players {
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

	// Neue Runde konfigurieren

	// explodierte Bomben zurücksetzen
	for i := range fields {
		for j := range fields[i] {
			f := fields[i][j]
			f.explodes = false
		}
	}

	// TODO: Specials respawn implementieren

	// nächste Runde setzen

	if round.id == len(manager.currentGame.rounds) {
		manager.currentGame.finished = true

		if manager.currentGame.id == len(manager.games) {
			manager.finishGame()
		} else {
			manager.nextGame()
		}

	} else {
		manager.currentGame.currentRound = manager.currentGame.rounds[round.id]

	}

	manager.broadcastWaiting()

	manager.timeout()

}

func (manager *Manager) nextGame() {
	previousGame := manager.currentGame
	manager.currentGame = manager.games[manager.currentGame.id]
	manager.currentGame.started = true
	manager.currentGame.gameMap = previousGame.gameMap // TODO: oder neue Map erzeugen
	manager.currentGame.players = previousGame.players
	manager.currentGame.currentRound = manager.currentGame.rounds[0]
}

func (manager *Manager) broadcastWaiting() {
	log.Println("waiting for commands")
	for _, p := range manager.currentGame.players {
		conn := manager.playersConn[p.id]
		conn.Write([]byte("wfyc: waiting for your command\n"))
	}
}

func (manager *Manager) broadcastGamestate() {
	for _, p := range manager.currentGame.players {
		manager.sendGameStateToPlayer(p)
	}

	log.Println(manager.GameState(manager.currentGame.gameMap.toStringForServer()))

	manager.specChannel <- manager.GameState(manager.currentGame.gameMap.toStringForServer())
}

func (manager *Manager) sendGameStateToPlayer(p *Player) {
	conn := manager.playersConn[p.id]
	conn.Write([]byte(buildHeader(manager.GameState(manager.currentGame.gameMap.toString())))) // ???
	conn.Write([]byte(manager.GameState(manager.currentGame.gameMap.toString())))
	conn.Write([]byte(p.msg))
	p.msg = ""
}

// Gibt für einen gegebenen Spieler und ein Ziel das entsprechende Feld zurück.
func (manager *Manager) destinationField(player *Player, destination []string) *Field {

	var distance int
	var direction string
	var err error

	if len(destination) < 3 {
		if (len(destination) == 2) {
			player.msg += "E: incomplete bomb command.\n"
		}
		return player.currentField
	} else {
		distance, err = strconv.Atoi(destination[1])
		if (err == nil) {
			direction = destination[2]
		} else {
			player.msg += "E: invalid distance " + destination[1] + ".\n"
			return player.currentField
		}
	}

	// prüfen ob Richtung gültig ist
	validDirections := "wasd"
	if strings.Contains(validDirections, direction) == false {
	        player.msg += "E: invalid direction " + direction + ". Must be in [wasd].\n"
		return player.currentField
	}

	pRow := player.currentField.row
	pCol := player.currentField.column

	var destinationField *Field

	// Prüfen ob der Player weiter werfen will als er darf (unmoeglich so lange distanz einstellig und throwrange immer 9).
	if distance > player.throwrange {
		distance = player.throwrange
	}
	
	// The direction-deltas while looking for a field to drop the bomb.
	var dx int
	var dy int
	
	switch direction {
	// Norden
	case "w":
		dx = 0
		dy = -1
	// Osten
	case "d":
		dx = +1
		dy = 0
	// Süden
	case "s":
		dx = 0
		dy = +1
	// Westen
	case "a":
		dx = -1
		dy = 0
	}
	
        for (distance > 0) && manager.currentGame.gameMap.isBombable(pRow + dy, pCol + dx) {
                distance -= 1
                pRow += dy
                pCol += dx
        }
        destinationField = manager.currentGame.gameMap.getField(pRow, pCol)

	return destinationField
}

func (manager *Manager) playerQuit(player *Player) {

	conn := manager.playersConn[player.id]
	conn.Write([]byte("good bye\n"))

	delete(manager.playersConn, player.id)
	manager.currentGame.removePlayer(player)

	log.Printf("Player %s has left the game.\n", player.id)

	conn.Close()
}

func (manager *Manager) finishGame() {
	log.Println("Game is over.")
	// TODO: Spielergebnis senden

	for _, conn := range manager.playersConn {
		conn.Write([]byte("game is over\n"))
		conn.Close()
	}

	os.Exit(0)
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
