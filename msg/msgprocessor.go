package msg

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"io"
)

type Processor struct {
	reader io.Reader
}

// Make new processor that will read the message from their
// binary format from reader and send them on sync
func NewProcessor(reader io.Reader, sync chan<- *Msg) *Processor {
	p := new(Processor)
	p.reader = bufio.NewReader(reader)

	return p
}

func (p *Processor) ReadMSG() (msg *Msg, err error) {
	return ReadMSG(p.reader)
}

func ReadMSG(reader io.Reader) (msg *Msg, err error) {
	msg = new(Msg)
	var buf []byte
	buf = make ([]byte, 4)
	_, err = io.ReadFull(reader, buf)
	if err != nil {
		return nil, err
	}

	msg.size = binary.BigEndian.Uint32(buf)
	buf = make ([]byte, msg.size + 1)
	_, err = io.ReadFull(reader, buf)
	if err != nil {
		return nil, err
	}
	msg.flag= buf[0]
	msg.payload = bytes.NewBuffer(buf[1:])
	return msg, nil
}
