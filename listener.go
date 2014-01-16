package gomahawk

import (
	"errors"
	"log"
	"net"
)

type tcpListener struct {
	addr     *net.TCPAddr
	callback func(*net.TCPConn) error
	listener *net.TCPListener
}

func newTCPListener(ip net.IP, port int, callback func(*net.TCPConn) error) (result *tcpListener, err error) {
	result = &tcpListener{
		addr:     &net.TCPAddr{IP: ip, Port: port},
		callback: callback,
	}
	result.listener, err = net.ListenTCP("tcp", result.addr)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (t tcpListener) Start() error {

	go func() {
		for {
			conn, err := t.listener.AcceptTCP()
			if err != nil {
				log.Println("error while accepting tcp with", t.listener, ":", err)
				if nerr, ok := err.(net.Error); ok && nerr.Temporary() {
					continue
				}

				break
			}

			t.callback(conn)
		}
	}()

	return nil
}

func (t tcpListener) Close() error {
	if t.listener == nil {
		return errors.New("Listener is either stopped or has never started")
	}

	return t.listener.Close()
}
