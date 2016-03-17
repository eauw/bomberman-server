package main

import (
	"bomberman-server/gamemanager"
	"bomberman-server/helper"
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	maximumPlayers = 2
)

var httpServer *HTTPServer
var httpChannel chan string
var mainChannel chan string

var mutex = &sync.Mutex{}

func main() {
	// create main channel
	mainChannel = make(chan string)
	go handleMainChannel()

	// handle command line arguments
	handleArgs()

	tcpPort := 5000
	fmt.Println("\n\n\n\nLaunching tcp server...")

	// listen on all interfaces
	ln, _ := net.Listen("tcp", fmt.Sprintf(":%d", tcpPort))
	fmt.Printf("Listening tcp on port %d\n", tcpPort)

	// create game
	gameManager := gamemanager.NewManager()
	gameManager.SetMainChannel(mainChannel)

	for {
		// accept connection on port
		conn, _ := ln.Accept()
		if conn != nil {
			go newClientConnected(conn, gameManager)
		}
	}
}

// check commandline arguments on program start
func handleArgs() {
	// only check if there is an parameter given
	if len(os.Args) > 1 {
		// ignore first parameter because its the programs name
		for i, v := range os.Args {
			if i > 0 {
				switch v {

				default:
					fmt.Println("invalid commandline parameter")
					os.Exit(0)
					break

				// show help
				case "-h":
					// show help
					helpString := "\nBomberman-Server is a game server for MICA 2016.\n\n"
					helpString += "Commands:\n\n"
					helpString += "\t-w\tstarts with http server\n\n"
					fmt.Print(helpString)
					os.Exit(0)
					break

				// start with http server
				case "-w":
					fmt.Println("Launching http server...")
					httpServer = NewHTTPServer()
					httpChannel = httpServer.channel
					httpServer.mainChannel = mainChannel
					// httpServer.game = game
					go httpServer.start()
					fmt.Printf("Listening http on port %s\n", httpServer.port)
					break

				}
			}
		}
	}
}

// called as goroutine
func handleMainChannel() {
	for {
		var x = <-mainChannel
		fmt.Print(x)
	}
}

// called as goroutine
func newClientConnected(conn net.Conn, gameManager *gamemanager.Manager) {
	fmt.Printf("\nclient %s connected\n", conn.RemoteAddr())
	conn.Write([]byte("Successfully connected to Bomberman-Server\n"))
	conn.Write([]byte("Enter quit or exit to disconnect.\n"))

	// get clients ip
	clientIP := helper.IpFromAddr(conn)

	mutex.Lock()
	newPlayer := gameManager.PlayerConnected(clientIP, conn)
	mutex.Unlock()

	if gameManager.PlayersCount() >= maximumPlayers {
		gameManager.Start()
	}

	conn.Write([]byte("Your ID: "))
	conn.Write([]byte(newPlayer.GetID()))
	conn.Write([]byte("\n"))

	// run loop forever (or until ctrl-c)
	for {
		messageBytes, _, err := bufio.NewReader(conn).ReadLine()
		if err == nil {
			messageString := string(messageBytes)

			// output message received
			fmt.Println("----------------")
			timeStamp := time.Now()
			fmt.Println(timeStamp)

			mainChannel <- fmt.Sprintf("Message from client: %s\n", clientIP)
			mainChannel <- fmt.Sprintf("Message Received:%s\n", messageString)

			// if message is "quit" server will close connection
			if messageString == "quit" {
				conn.Close()
				return
			} else {
				mutex.Lock()
				gameManager.MessageReceived(messageString, newPlayer)
				mutex.Unlock()
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

// handles the message sent by a client. returns a converted message and true. if message is "quit" or "exit" it returns false.
func handleMessage(message string) (string, bool) {
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
		return "", false

	case "exit":
		return "", false

	case "game state":
		// mainChannel <- game.gameMap.toString()
		break

	case "show players":
		// printMessage = game.printPlayers()
		break

	default:
		fmt.Printf("invalid command: %d", message)
		break
	}

	printMessage += "\n"
	mainChannel <- fmt.Sprintf(printMessage)

	return "", true
}
