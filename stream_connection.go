package gomahawk

import (
	"github.com/MStoykov/gomahawk/msg"
	"log"
)

type streamConn struct {
	secondaryConnection
	id     int64
	stream chan []byte
}

func openNewStreamConnection(id int64, c *connection, parent *connection, offerMsg *msg.Msg) (*streamConn, error) {

	conn := new(streamConn)
	conn.connection = c
	conn.parent = parent
	conn.id = id
	conn.stream = make(chan []byte)
	go func() {
		for m := range conn.sync {
			if err := conn.handleMsg(m); err != nil {
				log.Println("error in stream handling", err)
			}
		}
	}()
	return conn, nil
}

func (s *streamConn) FileID() int64 {
	return s.id
}

func (s *streamConn) BlockSize() int {
	return 4096
}

func (s *streamConn) SeekToBlock(blockIndex int) error {
	return nil
}

func (s *streamConn) GetStream() (<-chan []byte, error) {
	return s.stream, nil
}

func (s *streamConn) handleMsg(m *msg.Msg) error {
	if m.IsRaw() {
		s.stream <- m.Payload()[4:]
		if !m.IsFragment() {
			close(s.stream)
			return s.Close()
		}
	}
	return nil
}
