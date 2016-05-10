package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/eauw/bomberman-server/gamemanager"
	"github.com/eauw/bomberman-server/helper"
	"log"
	"net"
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
var maxPlayers int
var height int
var width int
var timeout float64
var gamesCount int

var mutex *sync.Mutex

func init() {
	flag.IntVar(&maxPlayers, "players", 2, "set min. players")
	flag.IntVar(&rounds, "rounds", 200, "set max. rounds")
	flag.IntVar(&height, "height", 20, "set maps height")
	flag.IntVar(&width, "width", 20, "set maps width")
	flag.IntVar(&gamesCount, "games", 3, "set how many games will be played")
	flag.Float64Var(&timeout, "timeout", 0.5, "command timeout")
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
	// mutex.Lock()
	gameManager := gamemanager.NewManager()
	gameManager.Start(rounds, height, width, gamesCount, timeout)
	gameManager.SetMainChannel(mainChannel)
	gameManager.SetSpecChannel(specChannel)
	// mutex.Unlock()

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
	conn.Write([]byte("Enter q to disconnect.\n"))

	if gameManager.PlayersCount() < maxPlayers {
		// get clients ip
		clientIP := helper.IpFromAddr(conn)

		mutex.Lock()
		newPlayer := gameManager.PlayerConnected(clientIP, conn)
		mutex.Unlock()

		conn.Write([]byte("YourID:"))
		conn.Write([]byte(newPlayer.GetID()))
		conn.Write([]byte("\n"))
		conn.Write([]byte("YourName:"))
		conn.Write([]byte(newPlayer.GetName()))
		conn.Write([]byte("\n"))

		if gameManager.PlayersCount() == maxPlayers {
			timer := time.NewTimer(time.Second)
			go func() {
				<-timer.C
				mutex.Lock()
				gameManager.GameStart()
				mutex.Unlock()
			}()

		}

		// run loop forever (or until ctrl-c)
		for {
			messageBytes, _, err := bufio.NewReader(conn).ReadLine()
			if err == nil {
				messageString := string(messageBytes)

				// output message received
				// fmt.Println("----------------")
				// timeStamp := time.Now()
				// fmt.Println(timeStamp)

				// mainChannel <- fmt.Sprintf("Message from client: %s\n", clientIP)
				// mainChannel <- fmt.Sprintf("Message Received:%s\n", messageString)

				gameMessage := gamemanager.NewGameChannelMessage(messageString, newPlayer)
				gameManager.GetGameChannel() <- gameMessage

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
	} else {
		conn.Write([]byte("Sorry, the player limit is reached.\n"))
		conn.Close()
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
