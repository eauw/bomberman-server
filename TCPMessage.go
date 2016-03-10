package main

type TCPMessage struct {
	text     string
	senderIP string
}

func NewTCPMessage(text string, senderIP string) *TCPMessage {

	return &TCPMessage{
		text:     text,
		senderIP: senderIP,
	}
}
