package tcpmessage

type TCPMessage struct {
	Text     string
	SenderIP string
}

func NewTCPMessage(text string, senderIP string) *TCPMessage {

	return &TCPMessage{
		Text:     text,
		SenderIP: senderIP,
	}
}
