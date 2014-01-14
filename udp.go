package gomahawk

import (
	"bytes"
	"log"
	"net"
	"strconv"

	gouuid "github.com/nu7hatch/gouuid"
)

type advertPacket struct {
	advert   string
	port     string
	uuid     *gouuid.UUID
	hostname *string
}

var DefaultUDPPort int = 50210

// make new advertPacket
func newAdvertPacket(advert string, port string, hostname string) (*advertPacket, error) {
	uuid, err := gouuid.NewV5(gouuid.NamespaceURL, []byte(hostname))
	if err != nil {
		return nil, err
	}
	packet := advertPacket{
		advert,
		port,
		uuid,
		&hostname,
	}
	return &packet, nil
}

func (a *advertPacket) Bytes() []byte {
	msg := bytes.NewBufferString(a.advert)
	msg.WriteRune(':')
	msg.WriteString(a.port)
	msg.WriteRune(':')
	msg.WriteString(a.uuid.String())
	if a.hostname != nil {
		msg.WriteRune(':')
		msg.WriteString(*a.hostname)
	}
	return msg.Bytes()
}

// sends an Advert One time
func advert(localAddr *net.UDPAddr, uuid *gouuid.UUID) error {
	packet := advertPacket{
		"TOMAHAWKADVERT",
		strconv.Itoa(localAddr.Port),
		uuid,
		nil,
	}
	socket, err := net.DialUDP("udp4", localAddr, &net.UDPAddr{
		IP:   net.ParseIP("255.255.255.255"),
		Port: DefaultUDPPort,
	})
	if err != nil {
		return err
	}
	log.Println(socket.RemoteAddr())
	log.Println(socket.LocalAddr())

	size, err := socket.Write(packet.Bytes())
	if err != nil {
		return err
	}
	log.Printf("wrote %d bytes\n", size)
	return nil
}

