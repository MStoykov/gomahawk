package gomahawk

import (
	"net"
	"time"
)

type Gomahawk interface {
	// Human readable Name
	Name() string
	// A new Connection from the given address is requested.
	// If the returned value is true GomahawkServer will continue
	// with the connection otherwise the connection will be closed.
	ConnectionIsRequested(net.Addr) bool
	// Found peer at the given address.
	// Name could be empty.
	// If the returned value is true GomahawkServer will try to connect to it.
	NewTomahawkFound(addr net.Addr, name string) bool
	// The given Tomahawk made DBConnection(by our request).
	//
	NewDBConnection(Tomahawk, DBConnection) error

	// The given Tomahawk has requested a DBConnection. The returned DBConnection will get a call to FetchOps.
	NewDBConnectionRequested(Tomahawk) (DBConnection, error)

	// The given Tomahawk has requested a StreamConnection for file with the given uuid
	// if the returned connection isn't nil it will be used to make the transactions
	NewStreamConnectionRequested(t Tomahawk, uuid string) (StreamConnection, error)

	// The given Tomahawk has opened StreamConnection to us.
	//
	// This StreamConnection is used for us to request blocks from them and is answer to
	// previous call to Tomahawk.RequestStraamConnection
	NewStreamConnection(Tomahawk, StreamConnection) error
}

// GomahawkServer
type GomahawkServer interface {
	// say that the instance needs to listen to the given port and ip.
	// error is being returned if that is not possible
	ListenTo(ip net.IP, port string) error
	// start the listening and advertising
	Start() error
	// says to advertise this instance
	// (this is regardless of advertisement period set by AdvertEvery)
	AdvertNow() error
	// sets a period of time that the advert will be sent
	// by default never
	AdvertEvery(period time.Duration)
	// retuns the name
	Name() string
	// returns the currently connected tomahawks
	GetTomahawks() []Tomahawk
}

// returns new instance of Gomahawk
func NewGomahawkServer(t Gomahawk) (result GomahawkServer, err error) {
	return
}
