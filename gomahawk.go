package gomahawk

import (
	"net"
	"time"

)

// Gomahawk 
type Gomahawk interface {
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
func NewGomahawk(t Tomahawk) (result Gomahawk, err error) {
	return
}

