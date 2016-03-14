package gamemanager

import (
	"fmt"
	"log"
	"sync"
)

type Game struct {
	channel     chan *GameChannelMessage
	mainChannel chan string
	gameMap     *GameMap
	players     map[string]*Player
}

func NewGame() *Game {
	ch := make(chan *GameChannelMessage)
	gm := NewGameMap(10)

	newGame := &Game{
		channel: ch,
		gameMap: gm,
		players: make(map[string]*Player),
	}

	gm.game = newGame

	return newGame
}

func (game *Game) start() {
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
	player.currentField = game.gameMap.getField(0, 0)
	game.gameMap.getField(0, 0).addPlayer(player) // .players = append(firstField.players, player)
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

// TODO: prüfen ob spieler map verlassen will

func (game *Game) PlayerMovesToLeft(player *Player) {
	mutex.Lock()
	defer mutex.Unlock()

	currentField := game.gameMap.fields[player.currentField.row][player.currentField.column]

	if currentField.column == 0 {
		return
	}

	nextField := game.gameMap.fields[currentField.row][currentField.column-1]
	nextField.addPlayer(player)
	currentField.removePlayer(player)

}

var mutex = &sync.Mutex{}

func (game *Game) PlayerMovesToRight(player *Player) {
	mutex.Lock()
	defer mutex.Unlock()

	currentField := game.gameMap.fields[player.currentField.row][player.currentField.column]

	if currentField.column == (game.gameMap.size - 1) {
		return
	}

	nextField := game.gameMap.fields[currentField.row][currentField.column+1]
	nextField.addPlayer(player)
	currentField.removePlayer(player)
}

func (game *Game) PlayerMovesToUp(player *Player) {
	mutex.Lock()
	defer mutex.Unlock()

	currentField := game.gameMap.fields[player.currentField.row][player.currentField.column]

	if currentField.row == 0 {
		return
	}

	nextField := game.gameMap.fields[currentField.row-1][currentField.column]
	nextField.addPlayer(player)
	currentField.removePlayer(player)

}

func (game *Game) PlayerMovesToDown(player *Player) {
	mutex.Lock()
	defer mutex.Unlock()

	currentField := game.gameMap.fields[player.currentField.row][player.currentField.column]

	if currentField.row == (game.gameMap.size - 1) {
		return
	}

	nextField := game.gameMap.fields[currentField.row+1][currentField.column]
	nextField.addPlayer(player)
	currentField.removePlayer(player)

}

// func (game *Game) placePlayers() {
// 	for p := range game.players {
//
// 	}
// }
