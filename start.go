package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

var game *Game
var httpServer *HTTPServer
var mainChannel chan string

func main() {
	// handle command line arguments
	handleArgs()

	mainChannel = make(chan string)
	go handleMainChannel()

	tcpPort := 5000
	fmt.Println("\n\n\n\nLaunching tcp server...")

	// listen on all interfaces
	ln, _ := net.Listen("tcp", fmt.Sprintf(":%d", tcpPort))
	fmt.Printf("Listening tcp on port %d\n", tcpPort)

	game = NewGame()
	game.gameMap.toString()

	go handleGameChannel()

	for {
		// accept connection on port
		conn, _ := ln.Accept()
		if conn != nil {
			go newClientConnected(conn, game)
		}
	}
}

func handleArgs() {
	if len(os.Args) > 1 {
		for _, v := range os.Args {
			switch v {
			case "-h":
				// show help
				helpString := "\nBomberman-Server is a game server for MICA 2016.\n\n"
				helpString += "Commands:\n\n"
				helpString += "\t-w\tstarts with http server\n\n"
				fmt.Print(helpString)
				os.Exit(0)
				break

			case "-w":
				fmt.Println("Launching http server...")
				httpServer = NewHTTPServer()
				httpServer.game = game
				go httpServer.start()
				fmt.Printf("Listening http on port %s\n", httpServer.port)
				go handleHTTPChannel()
			}

		}
	}
}

func handleMainChannel() {
	for {
		var x = <-mainChannel
		fmt.Print(x)
	}
}

func handleGameChannel() {
	for {
		var x = <-game.channel
		fmt.Printf("game channel: %s", x)
	}
}

func handleHTTPChannel() {
	for {
		var x = <-httpServer.channel
		fmt.Printf("httpServer: %s\n", x)
	}
}

func newClientConnected(conn net.Conn, game *Game) {
	fmt.Printf("\nclient %s connected\n", conn.RemoteAddr())
	conn.Write([]byte("Successfully connected to Bomberman-Server\nEnter quit or exit to disconnect.\n"))

	newPlayer := NewPlayer("New Player")
	game.addPlayer(newPlayer)

	fieldPlayers := game.gameMap.fields[0][0].players
	fieldPlayers = append(fieldPlayers, newPlayer)

	// run loop forever (or until ctrl-c)
	for {
		messageBytes, _, err := bufio.NewReader(conn).ReadLine()
		if err == nil {
			messageString := string(messageBytes)
			// output message received
			fmt.Println("----------------")
			timeStamp := time.Now()
			fmt.Println(timeStamp)
			fmt.Printf("Message from client: %s\n", conn.RemoteAddr())
			fmt.Printf("Message Received:%s\n", messageString)
			game.channel <- messageString

			// if message is "quit" server will close connection
			if handleMessage(messageString) == false {
				conn.Close()
				return
			}

			// sample process for string received
			newMessage := strings.ToUpper(messageString)
			// send new string back to client
			conn.Write([]byte(newMessage + "\n"))
		} else {
			fmt.Printf("Connection Error: %s\n", err)
			fmt.Println("Client disconnected.")
			conn.Close()
			return
		}
	}
}

func handleMessage(message string) bool {
	printMessage := ""

	switch message {
	case "a":
		printMessage = "go left"
		break

	case "d":
		printMessage = "go right"
		break

	case "w":
		printMessage = "go up"
		break

	case "s":
		printMessage = "go down"
		break

	case "quit":
		return false

	case "exit":
		return false

	case "game state":
		game.gameMap.toString()
		break

	default:
		fmt.Printf("no valid command: %d", len(message))
		break
	}

	fmt.Println(printMessage)

	return true
}
