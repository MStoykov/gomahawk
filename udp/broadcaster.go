package udp

import (
	"net"
)

type Broadcaster struct {
	addr   *net.UDPAddr
	socket *net.UDPConn
	msg    []byte
}

func NewBroadCaster(localIp net.IP, remotePort int, message []byte) (broadcaster *Broadcaster, err error) {
	broadcaster = new(Broadcaster)
	broadcaster.addr = &net.UDPAddr{IP: localIp, Port: 0}
	broadcaster.msg = message

	broadcaster.socket, err = net.DialUDP("udp", broadcaster.addr, &net.UDPAddr{
		IP:   net.IPv4bcast,
		Port: remotePort,
	})

	if err != nil {
		return nil, err
	}
	return broadcaster, nil
}

func (b *Broadcaster) Broadcast() (err error) {
	_, err = b.socket.Write(b.msg)

	return
}
