package main

import (
	"bomberman-server/gamemanager"
	"bomberman-server/helper"
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
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

var rounds = 20
var xSize = 20
var ySize = 20

var mutex *sync.Mutex

// check commandline arguments on program start
func handleArgs() {
	// only check if there is an parameter given
	if len(os.Args) > 1 {
		// ignore first parameter because its the programs name
		for i, v := range os.Args {
			if i > 0 {
				switch v {

				// show help
				case "-h":
					// show help
					showCommandlineHelp()
					break

				case "help":
					showCommandlineHelp()
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

				case "-r":
					rounds, _ = strconv.Atoi(os.Args[i+1])
					return
					break

				case "-s":
					xSize, _ = strconv.Atoi(os.Args[i+1])
					ySize, _ = strconv.Atoi(os.Args[i+2])
					return
					break

				default:
					fmt.Println("invalid commandline parameter")
					os.Exit(0)
					break

				}
			}
		}
	}
}

func main() {
	mutex = &sync.Mutex{}

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
	mutex.Lock()
	gameManager := gamemanager.NewManager()
	gameManager.Start(rounds, xSize, ySize)
	gameManager.SetMainChannel(mainChannel)
	mutex.Unlock()

	for {
		// accept connection on port
		conn, err := ln.Accept()
		if err != nil {
			log.Print(err)
		}

		if conn != nil {
			go newClientConnected(conn, gameManager)
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
		gameManager.GameStart()
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
