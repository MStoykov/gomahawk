package msg

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

const (
	RAW        byte = 1
	JSON            = 2
	FRAGMENT        = 4
	COMPRESSED      = 8
	DBOP            = 16
	PING            = 32
	RESERVED        = 64
	SETUP           = 128
)

// a generic Message between tomahawks
type Msg struct {
	payload []byte
	flag    byte
	size    uint32
}

// Create new Message with the given payload and flags
func NewMsg(payload []byte, flag byte) *Msg {
	return &Msg{
		payload,
		flag,
		uint32(len(payload)),
	}
}

func (t *Msg) IsRaw() bool {
	return t.flag&RAW == RAW
}

func (t *Msg) IsJSON() bool {
	return t.flag&JSON == JSON
}

func (t *Msg) IsFragment() bool {
	return t.flag&FRAGMENT == FRAGMENT
}

func (t *Msg) IsCompressed() bool {
	return t.flag&COMPRESSED == COMPRESSED
}

func (t *Msg) IsDBOP() bool {
	return t.flag&DBOP == DBOP
}

func (t *Msg) IsPing() bool {
	return t.flag&PING == PING
}

func (t *Msg) IsReserved() bool {
	return t.flag&RESERVED == RESERVED
}

func (t *Msg) IsSetup() bool {
	return t.flag&SETUP == SETUP
}

// Returns the payload of the message as byte array
func (t *Msg) Payload() []byte {
	return t.payload
}

// Uncompresses a compressed payload
func (t *Msg) Uncompress() {
	if t.IsCompressed() {
		t.flag ^= COMPRESSED
		t.payload = uncompress(t.payload)
	}
}

// Compresses an uncompressed payload
func (t *Msg) Compress() {
	if !t.IsCompressed() {
		t.flag ^= COMPRESSED
		t.payload = compress(t.payload)
	}
}

func (t *Msg) String() string {
	var payload []byte
	if t.IsCompressed() {
		payload = uncompress(t.payload)
	} else {
		payload = t.payload
	}

	return fmt.Sprintf("size: %d, flag: %s,\n payload[%s]", len(payload), t.flagToString(), payload)
}

func (t *Msg) flagToString() string {
	var buf bytes.Buffer
	if t.IsRaw() {
		buf.WriteString("Raw ")
	}
	if t.IsJSON() {
		buf.WriteString("JSON ")
	}
	if t.IsFragment() {
		buf.WriteString("Fragment ")
	}
	if t.IsCompressed() {
		buf.WriteString("Compressed ")
	}
	if t.IsDBOP() {
		buf.WriteString("DBOP ")
	}
	if t.IsPing() {
		buf.WriteString("Ping ")
	}
	if t.IsReserved() {
		buf.WriteString("!Reserved! ")
	}
	if t.IsSetup() {
		buf.WriteString("Setup ")
	}
	return buf.String()
}
func ReadMSG(reader io.Reader) (msg *Msg, err error) {
	buf := make([]byte, 5)
	_, err = io.ReadFull(reader, buf)
	if err != nil {
		return nil, err
	}

	msg = new(Msg)
	msg.size = binary.BigEndian.Uint32(buf[:4])
	msg.flag = buf[4]

	msg.payload = make([]byte, msg.size)
	_, err = io.ReadFull(reader, msg.payload)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

func (m *Msg) WriteTo(w io.Writer) (n int64, err error) {
	b := make([]byte, 5)
	b[4] = m.flag
	binary.BigEndian.PutUint32(b[:4], uint32(len(m.payload)))
	size, err := w.Write(b)
	n += int64(size)
	if err != nil {
		return n, err
	}
	size, err = w.Write(m.payload)
	n += int64(size)
	return n, err
}
