package msg

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
)

func compressBuffer(b *bytes.Buffer) *bytes.Buffer {
	var size uint32

	size = uint32(b.Len())

	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, size)

	w := zlib.NewWriter(buf)
	b.WriteTo(w)
	w.Close()

	return buf
}

func compress(b []byte) []byte {
	return compressBuffer(bytes.NewBuffer(b)).Bytes()
}

func uncompressBuffer(b *bytes.Buffer) *bytes.Buffer {
	var size uint32

	buf := bytes.NewBuffer(b.Next(4))
	binary.Read(buf, binary.BigEndian, &size)
	r, _ := zlib.NewReader(b)
	buf.Reset()
	buf.Grow(int(size))
	buf.ReadFrom(r)

	return buf

}

func uncompress(b []byte) []byte {
	return uncompressBuffer(bytes.NewBuffer(b)).Bytes()
}
