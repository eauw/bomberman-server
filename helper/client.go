package helper

import (
  "net"
  "strings"
)

func IpFromAddr(conn net.Conn) string {
  addr := conn.RemoteAddr()
  addrString := addr.String()
  clientIP := ""

  if strings.Contains(addrString,"::1") {
    clientIP = "localhost"
  } else {
    remAddr := strings.Split(addrString, ":")
  	clientIP = remAddr[0]
  }
  
  return clientIP
}
