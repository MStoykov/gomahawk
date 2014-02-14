// Implementation of the network protocol of tomahawk
//
// Basic Usage :
//
// 1.Implement Gomahawk
//
// 2.Make newInstance of GomahawkServer with NewGomahawkServer(GomahawkImpl)
//
// 3.GomahawkServer.Start()

package gomahawk

import (
	"errors"
	"net"
)

var (
	NotSupportedConnection = errors.New("Not Supported Connection")
)

type Gomahawk interface {
	// Human readable Name
	Name() string

	// A new connection from the given address is requested.
	// If the returned value is true GomahawkServer will continue
	// with the connection otherwise the connection will be closed.
	ConnectionIsRequested(net.Addr) bool

	// Found peer at the given address.
	// Name could be empty.
	// If the returned value is true GomahawkServer will try to connect to it.
	NewTomahawkFound(addr net.Addr, name string) bool

	// The given Tomahawk made DBConnection (by our request).
	//
	NewDBConnection(Tomahawk, DBConnection) (DBConnection, error)

	// The given Tomahawk has requested a DBConnection.
	// If the returned DBConneciton is nil the request will be denied.
	NewDBConnectionRequested(Tomahawk, DBConnection) (DBConnection, error)

	// The given Tomahawk has requested a StreamConnection for file with the given id
	// if the returned connection isn't nil it will be used to stream the file
	NewStreamConnectionRequested(t Tomahawk, uuid string) (StreamConnection, error)
}

// GomahawkServer
type GomahawkServer interface {
	// Listen to the given port and ip. The default port is 50210
	// error is being returned if that is not possible
	ListenTo(ip net.IP, port int) error
	// start the listening and advertising
	Start() error
	// says to advertise this instance
	// (this is regardless of advertisement period set by AdvertEvery)
	Advertise() error
	// Sets a period of time that the advert will be sent.
	//
	// By default 0 which means never
	AdvertiseEvery(seconds int)
	// retuns the name. This is the same as the Gomahawk.Name() for the Gomahawk instance
	// given to GomahawkServer
	Name() string
	// returns the currently connected tomahawks
	GetTomahawks() []Tomahawk
}
