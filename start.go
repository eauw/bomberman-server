package main

import (
  "net"
  "fmt"
  "bufio"
  "strings"
  "time"
)

var game *Game

func main() {
  port := 5000
  fmt.Println("\n\n\n\nLaunching tcp server...")

  // listen on all interfaces
  ln, _ := net.Listen("tcp", fmt.Sprintf(":%d", port))
  fmt.Printf("Listening tcp on port %d\n", port)

  game = NewGame()
  go handleChannel()

  for {
    // accept connection on port
    conn, _ := ln.Accept()
    if conn != nil {
      go newClientConnected(conn)
    }
  }
}

func handleChannel() {
  for {
    var x = <- game.channel
    fmt.Printf("channel: %s", x)
  }
}

func newClientConnected(conn net.Conn) {
  fmt.Printf("\nclient %s connected\n", conn.RemoteAddr())
  conn.Write([]byte("Successfully connected to Bomberman-Server\nEnter quit or exit to disconnect.\n"))

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

  default:
    fmt.Printf("no valid command: %d", len(message))
    break
  }

  fmt.Println(printMessage)

  return true
}
