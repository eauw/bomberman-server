package gamemanager

import (
	"bomberman-server/helper"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
)

type Game struct {
	channel     chan *GameChannelMessage
	mainChannel chan string
	gameMap     *GameMap
	players     map[string]*Player
	started     bool
}

func NewGame(xSize int, ySize int) *Game {
	ch := make(chan *GameChannelMessage)
	gm := NewGameMap(xSize, ySize)

	newGame := &Game{
		channel: ch,
		gameMap: gm,
		players: make(map[string]*Player),
		started: false,
	}

	gm.game = newGame

	return newGame
}

func (game *Game) start() {
	game.started = true
	game.pickRandomPlayer().isFox = 1
	game.placePlayers()

	go game.handleGameChannel()
}

// receives all information about the game
func (game *Game) handleGameChannel() {
	for {
		var gameChannelMessage = <-game.channel
		handleGameChannelMessage(gameChannelMessage)
	}
}

func (game *Game) placePlayers() {

	for _, p := range game.players {

		isWall := true

		for isWall {
			randomX := helper.RandomNumber(0, game.gameMap.xSize-1)
			randomY := helper.RandomNumber(0, game.gameMap.ySize-1)

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
		break

	case "move down":
		// gcm.GetPlayer().moveDown()
		break

	case "move left":
		// gcm.GetPlayer().moveLeft()
		break

	case "move right":
		// gcm.GetPlayer().moveRight()
		break

	}
}

var mutex = &sync.Mutex{}

func (game *Game) PlayerMovesToLeft(player *Player) {
	mutex.Lock()
	defer mutex.Unlock()

	if player.isParalyzed {
		return
	}

	currentField := player.currentField

	// prüfen ob der Spieler sich am Spielfeldrand befindet
	if currentField.column == 0 {
		return
	}

	nextField := game.gameMap.fields[currentField.row][currentField.column-1]

	// prüfen ob der Spieler versucht gegen eine Wand zu laufen
	if nextField.wall != nil {
		return
	}

	// prüfen ob das nächste Feld ein Special hat
	if nextField.special != nil {
		player.applySpecial(nextField.special)
		nextField.special = nil
	}

	nextField.addPlayer(player)
	currentField.removePlayer(player)

}

func (game *Game) PlayerMovesToRight(player *Player) {
	mutex.Lock()
	defer mutex.Unlock()

	if player.isParalyzed {
		return
	}

	currentField := player.currentField

	// prüfen ob der Spieler sich am Spielfeldrand befindet
	if currentField.column == (game.gameMap.xSize - 1) {
		return
	}

	nextField := game.gameMap.fields[currentField.row][currentField.column+1]

	// prüfen ob der Spieler versucht gegen eine Wand zu laufen
	if nextField.wall != nil {
		return
	}

	// prüfen ob das nächste Feld ein Special hat
	if nextField.special != nil {
		player.applySpecial(nextField.special)
		nextField.special = nil
	}

	nextField.addPlayer(player)
	currentField.removePlayer(player)
}

func (game *Game) PlayerMovesToUp(player *Player) {
	mutex.Lock()
	defer mutex.Unlock()

	if player.isParalyzed {
		return
	}

	currentField := player.currentField

	// prüfen ob der Spieler sich am Spielfeldrand befindet
	if currentField.row == 0 {
		return
	}

	nextField := game.gameMap.fields[currentField.row-1][currentField.column]

	// prüfen ob der Spieler versucht gegen eine Wand zu laufen
	if nextField.wall != nil {
		return
	}

	if nextField.special != nil {
		player.applySpecial(nextField.special)
		nextField.special = nil
	}

	// prüfen ob das nächste Feld ein Special hat
	nextField.addPlayer(player)
	currentField.removePlayer(player)

}

func (game *Game) PlayerMovesToDown(player *Player) {
	mutex.Lock()
	defer mutex.Unlock()

	if player.isParalyzed {
		return
	}

	currentField := player.currentField

	// prüfen ob der Spieler sich am Spielfeldrand befindet
	if currentField.row == (game.gameMap.ySize - 1) {
		return
	}

	nextField := game.gameMap.fields[currentField.row+1][currentField.column]

	// prüfen ob der Spieler versucht gegen eine Wand zu laufen
	if nextField.wall != nil {
		return
	}

	if nextField.special != nil {
		player.applySpecial(nextField.special)
		nextField.special = nil
	}

	// prüfen ob das nächste Feld ein Special hat
	nextField.addPlayer(player)
	currentField.removePlayer(player)

}

func (game *Game) PlayerPlacesBomb(player *Player, destinationField *Field) {
	mutex.Lock()
	defer mutex.Unlock()

	// bomb := destinationField.addNewBomb(player)
	bomb := player.getAvailableBomb()
	destinationField.addBomb(bomb)
	bomb.field = destinationField
	game.gameMap.addBomb(bomb)

	bomb.isPlaced = true
}

func (game *Game) ExplodeBombs() {
	game.gameMap.bombs = []*Bomb{}

	//var fields []*Field

	for i := range game.gameMap.fields {
		for j := range game.gameMap.fields[i] {
			if len(game.gameMap.fields[i][j].bombs) > 0 {
				for ib := range game.gameMap.fields[i][j].bombs {
					game.gameMap.fields[i][j].bombs[ib].explode(game.gameMap)
				}
			}

			game.gameMap.fields[i][j].bombs = []*Bomb{}

		}
	}
}

func (game *Game) ExplodePlayersBombs(player *Player) {
	for i := range player.bombs {
		if player.bombs[i].isPlaced {
			player.bombs[i].explode(game.gameMap)
		}
	}
}

// func (game *Game) placePlayers() {
// 	for p := range game.players {
//
// 	}
// }
