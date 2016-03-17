package gamemanager

import (
	"fmt"
	"log"
	"sync"
)

type Game struct {
	channel      chan *GameChannelMessage
	mainChannel  chan string
	gameMap      *GameMap
	players      map[string]*Player
	rounds       int
	currentRound int
	started      bool
}

func NewGame() *Game {
	ch := make(chan *GameChannelMessage)
	gm := NewGameMap(10)

	newGame := &Game{
		channel:      ch,
		gameMap:      gm,
		players:      make(map[string]*Player),
		currentRound: 1,
		rounds:       20,
		started:      false,
	}

	gm.game = newGame

	return newGame
}

func (game *Game) start() {
	game.started = true
	go game.handleGameChannel()
}

// receives all information about the game
func (game *Game) handleGameChannel() {
	for {
		var gameChannelMessage = <-game.channel
		handleGameChannelMessage(gameChannelMessage)
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

	if nextField.special != nil {
		player.hasSpecial = true
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
	if currentField.column == (game.gameMap.size - 1) {
		return
	}

	nextField := game.gameMap.fields[currentField.row][currentField.column+1]

	// prüfen ob der Spieler versucht gegen eine Wand zu laufen
	if nextField.wall != nil {
		return
	}

	if nextField.special != nil {
		player.hasSpecial = true
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
		player.hasSpecial = true
		nextField.special = nil
	}

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
	if currentField.row == (game.gameMap.size - 1) {
		return
	}

	nextField := game.gameMap.fields[currentField.row+1][currentField.column]

	// prüfen ob der Spieler versucht gegen eine Wand zu laufen
	if nextField.wall != nil {
		return
	}

	if nextField.special != nil {
		player.hasSpecial = true
		nextField.special = nil
	}

	nextField.addPlayer(player)
	currentField.removePlayer(player)

}

func (game *Game) PlayerPlacesBomb(player *Player) {
	mutex.Lock()
	defer mutex.Unlock()

	currentField := player.currentField

	bomb := currentField.addNewBomb(player)
	game.gameMap.addBomb(bomb)
}

func (game *Game) ExplodeBomb() {
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

// func (game *Game) placePlayers() {
// 	for p := range game.players {
//
// 	}
// }
