package msg

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"io"
	"log"
)

type Processor struct {
	sync    chan<- *Msg
	current *Msg
	reader  *bufio.Reader
}

// Make new processor that will read the message from their 
// binary format from reader and send them on sync
func NewProcessor(reader io.Reader, sync chan<- *Msg) *Processor {
	p := Processor{
		sync: sync,
	}
	p.reader = bufio.NewReader(reader)

	go func() {
		defer close(sync)
		for {
			err := p.process()
			if err != nil {
				log.Println("error in process", err)
				return
			}
		}
	}()
	return &p
}

func readExactly(r *bufio.Reader, size uint32) (*bytes.Buffer, error) {
	var buf *bytes.Buffer
	buf = new(bytes.Buffer)

	buf.Grow(int(size)) // we don't actually want to reallocate every 10 calls

	var b [512]byte
	for newSize := size; newSize != 0; newSize = size - uint32(buf.Len()) {
		if newSize > 512 {
			newSize = 512
		}
		newSize, err := r.Read(b[:newSize])
		if err != nil {
			return nil, err
		}
		buf.Write(b[:newSize])
	}

	return buf, nil
}

func (p *Processor) process() error {
	buf, err := readExactly(p.reader, 4)

	p.current = new(Msg)
	if err != nil {
		if err != io.EOF {
			log.Println("Didn't manage to read size")
		} else {
			log.Println("EOF")
		}
		return err
	}

	err = binary.Read(buf, binary.BigEndian, &p.current.size)
	if err != nil {
		panic("error while parsing the size of a packet")
	}
	p.current.payload, err = readExactly(p.reader, p.current.size+1)
	if err != nil {
		if err != io.EOF {
			log.Println("Didn't manage to read size")
		} else {
			log.Println("EOF")
		}
		return err
	}
	p.current.flag, err = p.current.payload.ReadByte()
	if err != nil {
		if err != io.EOF {
			log.Println("Didn't manage to read the one byte flag")
		}
		return err
	}
	p.sync <- p.current
	return nil
}
