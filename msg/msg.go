package msg

import (
	"bytes"
	"encoding/binary"
	"fmt"
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

// returns the byte representation of the message
func (t *Msg) Bytes() []byte {
	b := make([]byte, len(t.payload)+5)
	binary.BigEndian.PutUint32(b[:4], uint32(len(t.payload)))

	b[4] = t.flag
	copy(b[5:], t.payload)
	return b
}

// Uncompresses a compressed payload
func (t *Msg) Uncompress() {
	if !t.IsCompressed() {
		return
	}

	t.flag ^= COMPRESSED

	t.payload = uncompress(t.payload)

}

// Compresses an uncompressed payload
func (t *Msg) Compress() {
	if t.IsCompressed() {
		return
	}

	t.flag ^= COMPRESSED

	t.payload = compress(t.payload)

}

// decodes message from it's binary representation
func ParseMsg(b []byte) (*Msg, error) {
	var size uint32
	result := new(Msg)
	size = binary.BigEndian.Uint32(b[:4])
	result.flag = b[4]
	payload := b[5:]
	result.size = size
	result.payload = payload
	return result, nil
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
