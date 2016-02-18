package main

import "net"
import "fmt"
import "bufio"
import "strings" // only needed below for sample processing

func main() {
  port := 5000
  fmt.Println("Launching server...")

  // listen on all interfaces
  ln, _ := net.Listen("tcp", fmt.Sprintf(":%d",port))
  fmt.Printf("Listening on port %d\n", port)

  // accept connection on port
  conn, _ := ln.Accept()
  fmt.Printf("client %s connected\n", conn.RemoteAddr())

  // run loop forever (or until ctrl-c)
  for {
    message, _, _ := bufio.NewReader(conn).ReadLine()
    messageString := string(message)
    // output message received
    fmt.Printf("Message Received:%s\n", messageString)

    var printMessage = ""

    switch messageString {
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
      return

    default:
      fmt.Printf("no valid command")
      break
    }

    fmt.Println(printMessage)

    // sample process for string received
    newMessage := strings.ToUpper(messageString)
    // send new string back to client
    conn.Write([]byte(newMessage + "\n"))
  }
}
