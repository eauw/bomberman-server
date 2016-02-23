package main

import "net"
import "fmt"
import "bufio"
import "strings"

func main() {
  port := 5000
  fmt.Println("Launching tcp server...")

  // listen on all interfaces
  ln, _ := net.Listen("tcp", fmt.Sprintf(":%d",port))
  fmt.Printf("Listening tcp on port %d\n", port)

  // accept connection on port
  conn, _ := ln.Accept()
  fmt.Printf("client %s connected\n", conn.RemoteAddr())

  // run loop forever (or until ctrl-c)
  for {
    messageBytes, _, _ := bufio.NewReader(conn).ReadLine()
    messageString := string(messageBytes)
    // output message received
    fmt.Printf("Message Received:%s\n", messageString)

    // if message is "quit" server will stop
    if handleMessage(messageString, conn) == false {
      return
    }

    // sample process for string received
    newMessage := strings.ToUpper(messageString)
    // send new string back to client
    conn.Write([]byte(newMessage + "\n"))
  }
}

func handleMessage(message string, conn net.Conn) bool {
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
    conn.Close()
    return false

  default:
    fmt.Printf("no valid command")
    break
  }

  fmt.Println(printMessage)

  return true
}
