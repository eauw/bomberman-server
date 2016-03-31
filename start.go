package main

import (
	"bomberman-server/gamemanager"
	"bomberman-server/helper"
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

var httpServer *HTTPServer
var httpServerBool bool
var httpChannel chan string
var mainChannel chan string

var specChannel chan string

var spectators []net.Conn = make([]net.Conn, 0)

var rounds int
var minPlayers int
var xSize int
var ySize int

var mutex *sync.Mutex

func init() {
	flag.IntVar(&minPlayers, "p", 2, "set min. players")
	flag.IntVar(&rounds, "r", 1000, "set max. rounds")
	flag.IntVar(&xSize, "x", 20, "set maps x size")
	flag.IntVar(&ySize, "y", 20, "set maps y size")
	flag.BoolVar(&httpServerBool, "w", false, "start http server")
	flag.Parse()
}

func startHttpServer() {
	fmt.Println("Launching http server...")
	httpServer = NewHTTPServer()
	httpChannel = httpServer.channel
	httpServer.mainChannel = mainChannel
	// httpServer.game = game
	go httpServer.start()
	fmt.Printf("Listening http on port %s\n", httpServer.port)
}

func main() {
	mutex = &sync.Mutex{}

	// handle command line arguments
	if httpServerBool {
		startHttpServer()
	}

	// create main channel
	mainChannel = make(chan string)
	go handleMainChannel()

	specChannel = make(chan string)
	go handleSpecChannel()

	tcpGamePort := 5000
	tcpSpecPort := 5001
	fmt.Println("\n\n\n\nLaunching tcp servers...")

	// listen on all interfaces
	gameListener, _ := net.Listen("tcp", fmt.Sprintf(":%d", tcpGamePort))
	fmt.Printf("tcp game port: %d\n", tcpGamePort)

	specListener, _ := net.Listen("tcp", fmt.Sprintf(":%d", tcpSpecPort))
	fmt.Printf("tcp spectator port: %d\n", tcpSpecPort)

	// create game
	mutex.Lock()
	gameManager := gamemanager.NewManager()
	gameManager.Start(rounds, xSize, ySize)
	gameManager.SetMainChannel(mainChannel)
	gameManager.SetSpecChannel(specChannel)
	mutex.Unlock()

	go handleSpecListener(specListener)

	for {
		// accept connection on port
		gameConn, gameConnErr := gameListener.Accept()

		if gameConnErr != nil {
			log.Print(gameConnErr)
		}

		if gameConn != nil {
			go newClientConnected(gameConn, gameManager)
		}
	}
}

func handleSpecListener(ln net.Listener) {
	for {
		// accept connection on port
		conn, specConnErr := ln.Accept()

		if specConnErr != nil {
			log.Print(specConnErr)
		}

		if conn != nil {
			mutex.Lock()
			spectators = append(spectators, conn)
			mutex.Unlock()
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

	if gameManager.PlayersCount() >= minPlayers {
		gameManager.GameStart()
	}

	conn.Write([]byte("Your ID: "))
	conn.Write([]byte(newPlayer.GetID()))
	conn.Write([]byte("\n"))
	conn.Write([]byte("Your Name: "))
	conn.Write([]byte(newPlayer.GetName()))
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

			mutex.Lock()
			gameManager.MessageReceived(messageString, newPlayer)
			mutex.Unlock()

			// sample process for string received
			// newMessage := strings.ToUpper(messageString)
			// send new string back to client
			conn.Write([]byte(messageString + "\n"))
		} else {
			if strings.Contains(err.Error(), "use of closed network connection") {
				fmt.Printf("Client %s disconnected.\n", newPlayer.GetID())
			} else {
				fmt.Printf("Connection Error: %s\n", err)
				fmt.Println("Client disconnected.")
				conn.Close()
			}

			return
		}
	}
}

// called as goroutine
func handleSpecChannel() {
	for {
		var x = <-specChannel

		for _, conn := range spectators {
			conn.Write([]byte(x))
		}
	}

}

// called as goroutine
// func newSpectatorConnected(conn net.Conn) {
// 	conn.Write([]byte("Hello Spectator\n"))

// 	for {
// 		var x = <-specChannel

// 	}
// }

func showCommandlineHelp() {
	helpString := "\nBomberman-Server is a game server for MICA 2016.\n\n"
	helpString += "Commands:\n"
	helpString += "  -w                  starts with http server\n"
	helpString += fmt.Sprintf("  -r [int]            set number of rounds, default: %d\n", rounds)
	helpString += fmt.Sprintf("  -s [x int] [y int]  set map size, default: x %d y%d\n", xSize, ySize)
	helpString += "\n"
	fmt.Print(helpString)
	os.Exit(0)
}
