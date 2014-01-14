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

func newTCPListener(addr *net.TCPAddr, callback func(*net.TCPConn) error) (*tcpListener, error) {
	result := &tcpListener{addr, callback, nil}
	err := result.Start()
	if err != nil {
		result = nil
	}
	return result, err
}

func (t tcpListener) Start() error {
	listener, err := net.ListenTCP("tcp", t.addr)
	if err != nil {
		return err
	}

	go func() {
		log.Printf("Start listening on %s", t.addr)
		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				log.Println("error while accepting tcp with", listener, ":", err)
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
