package gomahawk

import (
	"errors"
	"net"
	"sync"
)

type connectionManager struct {
	m map[*connection]*net.TCPConn
	sync.Mutex
}

func newConnectionManager() *connectionManager {
	cm := new(connectionManager)
	cm.m = make(map[*connection]*net.TCPConn)
	return cm

}

func (cm *connectionManager) copyConnection(conn *connection) (*connection, error) {
	tcpConn, ok := cm.m[conn]
	if !ok {
		return nil, errors.New("not registered connection tried to be copied")
	}

	tcpLAddr, err := net.ResolveTCPAddr("tcp", tcpConn.LocalAddr().String())
	if err != nil {
		return nil, err
	}
	tcpLAddr.Port = 0
	tcpRAddr, err := net.ResolveTCPAddr("tcp", tcpConn.RemoteAddr().String())
	if err != nil {
		return nil, err
	}
	tcpRAddr.Port = 50210 // hardcoded

	newTCPConn, err := net.DialTCP("tcp", tcpLAddr, tcpRAddr)
	if err != nil {
		return nil, err
	}

	return cm.newConnection(newTCPConn)
}

func (cm *connectionManager) newConnection(tcpConn *net.TCPConn) (*connection, error) {
	cm.Lock()
	defer cm.Unlock()
	c := new(connection)
	c.conn = tcpConn
	c.setupProcessor()
	cm.m[c] = tcpConn
	return c, nil
}
