package gomahawk

import (
	"errors"
	"net"
	"sync"
)

type connectionManager struct {
	connections map[*connection]*net.TCPConn
	sync.Mutex
}

func newConnectionManager() *connectionManager {
	cm := new(connectionManager)
	cm.connections = make(map[*connection]*net.TCPConn)
	return cm

}

func (cm *connectionManager) copyConnection(conn *connection) (*connection, error) {
	tcpConn, ok := cm.connections[conn]
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

	return cm.registerConnection(newTCPConn), nil
}

func (cm *connectionManager) registerConnection(tcpConn *net.TCPConn) (conn *connection) {
	cm.Lock()
	defer cm.Unlock()
	conn = new(connection)
	conn.conn = tcpConn
	cm.connections[conn] = tcpConn
	return conn
}
