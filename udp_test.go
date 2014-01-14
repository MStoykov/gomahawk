package gomahawk

import (
	"bytes"
	"testing"
)

func TestNewAdvertPacket(t *testing.T) {
	advert, err := newAdvertPacket("TOMAHAWKADVERT", "50120", "example.com")
	if err != nil {
		t.Fatal("got error while making advert", err)
	}
	b := advert.Bytes()
	expected := []byte("TOMAHAWKADVERT:50120:a5cf6e8e-4cfa-5f31-6804-6de6d1245e26:example.com")
	if !bytes.Equal(b, expected) {
		t.Errorf("expected: %s\n got     : %s", expected, b)
	}
}
