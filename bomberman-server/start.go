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

  // accept connection on port
  conn, _ := ln.Accept()

  // run loop forever (or until ctrl-c)
  for {
    // will listen for message to process ending in newline (\n)
    message, _ := bufio.NewReader(conn).ReadString('\n')
    messageString := strings.Replace(message,"\r\n","",-1)
    // output message received
    fmt.Printf("Message Received:%s\n", messageString)
    //fmt.Printf(">>% x<<\n", messageString)

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
    newMessage := strings.ToUpper(message)
    // send new string back to client
    conn.Write([]byte(newMessage + "\n"))
  }
}
