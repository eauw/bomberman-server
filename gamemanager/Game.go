package gamemanager

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/eauw/bomberman-server/helper"
)

type Game struct {
	id           int
	channel      chan *GameChannelMessage
	mainChannel  chan string
	gameMap      *GameMap
	players      map[string]*Player
	started      bool
	finished     bool
	rounds       []*Round
	currentRound *Round
}

var mutex = &sync.Mutex{}

func NewGame(height int, width int) *Game {
	ch := make(chan *GameChannelMessage)
	gm := NewGameMap(height, width)

	newGame := &Game{
		channel:  ch,
		gameMap:  gm,
		players:  make(map[string]*Player),
		started:  false,
		finished: false,
		rounds:   []*Round{},
	}

	return newGame
}

func (game *Game) start() {
	log.Println("game start")
	game.started = true
	// game.pickRandomPlayer().isFox = 1
	game.placePlayers()
}

func (game *Game) placePlayers() {

	for _, p := range game.players {

		isWall := true

		for isWall {
			randomX := helper.RandomNumber(0, game.gameMap.height-1)
			randomY := helper.RandomNumber(0, game.gameMap.width-1)

			field := game.gameMap.fields[randomX][randomY]

			if field.wall == nil {
				isWall = false
				p.currentField.removePlayer(p)
				field.addPlayer(p)
				p.currentField = field
			}
		}

	}
}

func (game *Game) GetPlayers() map[string]*Player {
	return game.players
}

func (game *Game) GetPlayersArray() []*Player {
	players := []*Player{}
	for _, v := range game.players {
		players = append(players, v)
	}
	return players
}

func (game *Game) addPlayer(player *Player) {
	player.currentField = game.gameMap.fields[0][0]
	game.gameMap.fields[0][0].addPlayer(player) // .players = append(firstField.players, player)
	game.players[player.id] = player
}

func (game *Game) removePlayer(player *Player) {
	delete(game.players, player.id)
}

func (game *Game) printPlayers() string {
	s := "Players\n"
	i := 0
	for _, v := range game.players {
		playerCount := i + 1
		i++
		s += fmt.Sprintf("%d. ID: %s | IP: %s | Name: %s \n", playerCount, v.id, v.ip, v.name)
	}

	return s
}

func (game *Game) getPlayerByIP(ip string) *Player {
	for _, p := range game.players {
		if p.ip == ip {
			return p
		}
	}

	return nil
}

func (game *Game) getPlayerByID(id string) *Player {
	return game.players[id]
}

func (game *Game) pickRandomPlayer() *Player {
	rand.Seed(time.Now().UnixNano())
	i := rand.Intn(len(game.players))
	pArr := game.GetPlayersArray()
	return pArr[i]
}

func handleGameChannelMessage(gcm *GameChannelMessage) {
	log.Printf("aaaa: %s", gcm.player.toString())
	switch gcm.GetText() {
	case "move up":
		// gcm.GetPlayer().moveUp()

	case "move down":
		// gcm.GetPlayer().moveDown()

	case "move left":
		// gcm.GetPlayer().moveLeft()

	case "move right":
		// gcm.GetPlayer().moveRight()

	}
}

func (game *Game) PlayerMovesToLeft(player *Player) {
	mutex.Lock()
	defer mutex.Unlock()

	if player.isParalyzed > 0 {
		return
	}

	currentField := player.currentField

	// prüfen ob der Spieler sich am Spielfeldrand befindet
	if currentField.column == 0 {
		return
	}

	nextField := game.gameMap.fields[currentField.row][currentField.column-1]

	// prüfen ob der Spieler versucht gegen eine Wand zu laufen
	if nextField.wall != nil || len(nextField.bombs) > 0 {
		return
	}

	// prüfen ob das nächste Feld ein Special hat
	if nextField.special != nil {
		player.applySpecial(nextField.special)
		nextField.special = nil
	}

	// prüfen ob der Fuchs auf dem nächsten Feld steht
	// nur wenn man selbst nicht der Fuchs ist
	// if player.isFox == 0 {
	// 	for _, p := range nextField.players {
	// 		if p.isFox > 0 {
	// 			p.isFox = 0
	// 			player.isFox += 1
	// 			game.teleportPlayer(player)
	// 			return
	// 		}
	// 	}
	// }

	nextField.addPlayer(player)
	currentField.removePlayer(player)

}

func (game *Game) PlayerMovesToRight(player *Player) {
	mutex.Lock()
	defer mutex.Unlock()

	if player.isParalyzed > 0 {
		return
	}

	currentField := player.currentField

	// prüfen ob der Spieler sich am Spielfeldrand befindet
	if currentField.column == (game.gameMap.height - 1) {
		return
	}

	nextField := game.gameMap.fields[currentField.row][currentField.column+1]

	// prüfen ob der Spieler versucht gegen eine Wand zu laufen
	if nextField.wall != nil || len(nextField.bombs) > 0 {
		return
	}

	// prüfen ob das nächste Feld ein Special hat
	if nextField.special != nil {
		player.applySpecial(nextField.special)
		nextField.special = nil
	}

	// prüfen ob der Fuchs auf dem nächsten Feld steht
	// nur wenn man selbst nicht der Fuchs ist
	// if player.isFox == 0 {
	// 	for _, p := range nextField.players {
	// 		if p.isFox > 0 {
	// 			p.isFox = 0
	// 			player.isFox += 1
	// 			game.teleportPlayer(player)
	// 			return
	// 		}
	// 	}
	// }

	nextField.addPlayer(player)
	currentField.removePlayer(player)
}

func (game *Game) PlayerMovesToUp(player *Player) {
	mutex.Lock()
	defer mutex.Unlock()

	if player.isParalyzed > 0 {
		return
	}

	currentField := player.currentField

	// prüfen ob der Spieler sich am Spielfeldrand befindet
	if currentField.row == 0 {
		return
	}

	nextField := game.gameMap.fields[currentField.row-1][currentField.column]

	// prüfen ob der Spieler versucht gegen eine Wand zu laufen
	if nextField.wall != nil || len(nextField.bombs) > 0 {
		return
	}

	// prüfen ob das nächste Feld ein Special hat
	if nextField.special != nil {
		player.applySpecial(nextField.special)
		nextField.special = nil
	}

	// prüfen ob der Fuchs auf dem nächsten Feld steht
	// nur wenn man selbst nicht der Fuchs ist
	// if player.isFox == 0 {
	// 	for _, p := range nextField.players {
	// 		if p.isFox > 0 {
	// 			p.isFox = 0
	// 			player.isFox += 1
	// 			game.teleportPlayer(player)
	// 			return
	// 		}
	// 	}
	// }

	nextField.addPlayer(player)
	currentField.removePlayer(player)

}

func (game *Game) PlayerMovesToDown(player *Player) {
	mutex.Lock()
	defer mutex.Unlock()

	if player.isParalyzed > 0 {
		return
	}

	currentField := player.currentField

	// prüfen ob der Spieler sich am Spielfeldrand befindet
	if currentField.row == (game.gameMap.width - 1) {
		return
	}

	nextField := game.gameMap.fields[currentField.row+1][currentField.column]

	// prüfen ob der Spieler versucht gegen eine Wand zu laufen
	if nextField.wall != nil || len(nextField.bombs) > 0 {
		return
	}

	// prüfen ob das nächste Feld ein Special hat
	if nextField.special != nil {
		player.applySpecial(nextField.special)
		nextField.special = nil
	}

	// prüfen ob der Fuchs auf dem nächsten Feld steht
	// nur wenn man selbst nicht der Fuchs ist
	// if player.isFox == 0 {
	// 	for _, p := range nextField.players {
	// 		if p.isFox > 0 {
	// 			p.isFox = 0
	// 			player.isFox += 1
	// 			game.teleportPlayer(player)
	// 			return
	// 		}
	// 	}
	// }

	nextField.addPlayer(player)
	currentField.removePlayer(player)

}

func (game *Game) PlayerPlacesBomb(player *Player, destinationField *Field) {
	mutex.Lock()
	defer mutex.Unlock()

	bomb := player.getAvailableBomb()
	if bomb != nil {
		destinationField.addBomb(bomb)
		bomb.field = destinationField
		game.gameMap.addBomb(bomb)

	}
}

// func (game *Game) ExplodeBombs() {
// 	game.gameMap.bombs = []*Bomb{}

// 	fields := game.gameMap.fields

// 	for i := range fields {
// 		for j := range fields[i] {
// 			if len(fields[i][j].bombs) > 0 {
// 				for ib := range fields[i][j].bombs {
// 					fields[i][j].bombs[ib].explode(game.gameMap)
// 				}
// 			}

// 			fields[i][j].bombs = []*Bomb{}

// 		}
// 	}
// }

// func (game *Game) ExplodePlayersBombs(player *Player) {
// 	for i := range player.bombs {
// 		if player.bombs[i].isPlaced {
// 			player.bombs[i].explode(game.gameMap)
// 		}
// 	}
// }

func (game *Game) teleportPlayer(player *Player) {

	gameMap := game.gameMap

	randomX := helper.RandomNumber(0, gameMap.height-1)
	randomY := helper.RandomNumber(0, gameMap.width-1)

	field := game.gameMap.fields[randomX][randomY]

	// den Spieler nur platzieren wenn das Feld keine Wand hat und kein anderer Spieler dort steht
	if field.wall != nil || len(field.players) > 0 {
		game.teleportPlayer(player)
	} else {
		player.currentField.removePlayer(player)
		field.addPlayer(player)
		player.currentField = field
	}

}
