package msg

// Makes new Ping Message
func MakePingMsg() *Msg {
	return &Msg{
		[]byte{},
		PING,
		0,
	}
}
