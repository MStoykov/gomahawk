package msg

import (
	"bytes"
)

// Makes new Ping Message
func MakePingMsg() *Msg {
	return &Msg{
		&bytes.Buffer{},
		PING,
		0,
	}
}
