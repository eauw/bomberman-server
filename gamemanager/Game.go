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

func (game *Game) PlayerMovesToLeft(player *Player) {
	mutex.Lock()
	defer mutex.Unlock()
	i := player.currentField.horizontalFieldCode
	j := player.currentField.verticalFieldCode
	game.gameMap.fields[i][j].removePlayer(player)
	game.gameMap.fields[i-1][j].addPlayer(player)
	player.currentField = game.gameMap.fields[i-1][j]
	log.Print(player.currentField.toString())

	// currentField := game.gameMap.fields[player.currentField.horizontalFieldCode][player.currentField.verticalFieldCode]
	// nextField := game.gameMap.fields[currentField.horizontalFieldCode-1][currentField.verticalFieldCode]
	// nextField.addPlayer(player)
	// currentField.removePlayer(player)
}

var mutex = &sync.Mutex{}

func (game *Game) PlayerMovesToRight(player *Player) {
	mutex.Lock()
	defer mutex.Unlock()
	log.Println(player)
	i := player.currentField.horizontalFieldCode
	j := player.currentField.verticalFieldCode
	log.Println(i)
	log.Println(j)

	oldField := game.gameMap.getField(i, j)
	log.Println(oldField)
	oldField.removePlayer(player)
	log.Println(oldField)
	// game.gameMap.fields[i][j].removePlayer(player)
	i += 1
	newField := game.gameMap.getField(i, j)
	log.Println(newField)
	newField.addPlayer(player)
	log.Println(newField)
	// game.gameMap.fields[i+1][j].addPlayer(player)
	player.currentField = newField
	log.Print(player.currentField.toString())

	// currentField := game.gameMap.fields[player.currentField.horizontalFieldCode][player.currentField.verticalFieldCode]
	// nextField := game.gameMap.fields[currentField.horizontalFieldCode+1][currentField.verticalFieldCode]
	// nextField.addPlayer(player)
	// currentField.removePlayer(player)
}

func (game *Game) PlayerMovesToUp(player *Player) {
	mutex.Lock()
	defer mutex.Unlock()
	i := player.currentField.horizontalFieldCode
	j := player.currentField.verticalFieldCode
	game.gameMap.fields[i][j].removePlayer(player)
	game.gameMap.fields[i][j-1].addPlayer(player)
	player.currentField = game.gameMap.fields[i][j-1]
	log.Print(player.currentField.toString())

	// currentField := game.gameMap.fields[player.currentField.horizontalFieldCode][player.currentField.verticalFieldCode]
	// nextField := game.gameMap.fields[currentField.horizontalFieldCode][currentField.verticalFieldCode-1]
	// nextField.addPlayer(player)
	// currentField.removePlayer(player)
}

func (game *Game) PlayerMovesToDown(player *Player) {
	mutex.Lock()
	defer mutex.Unlock()
	i := player.currentField.horizontalFieldCode
	j := player.currentField.verticalFieldCode
	game.gameMap.fields[i][j].removePlayer(player)
	game.gameMap.fields[i][j+1].addPlayer(player)
	player.currentField = game.gameMap.fields[i][j+1]
	log.Print(player.currentField.toString())

	// currentField := game.gameMap.fields[player.currentField.horizontalFieldCode][player.currentField.verticalFieldCode]
	// nextField := game.gameMap.fields[currentField.horizontalFieldCode][currentField.verticalFieldCode+1]
	// nextField.addPlayer(player)
	// currentField.removePlayer(player)
}

// func (game *Game) placePlayers() {
// 	for p := range game.players {
//
// 	}
// }
