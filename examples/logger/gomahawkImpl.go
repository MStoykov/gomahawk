package logger

import (
	"net"

	"github.com/MStoykov/gomahawk"
)

type gomahawkImpl struct {
	name string
}

func (g *gomahawkImpl) Name() string {
	return g.name
}
func (g *gomahawkImpl) ConnectionIsRequested(addr net.Addr) bool {
	return true
}
func (g *gomahawkImpl) NewTomahawkFound(addr net.Addr, name string) bool {
	return true
}
func (g *gomahawkImpl) NewDBConnection(t gomahawk.Tomahawk, db gomahawk.DBConnection) error {
	return nil
}
func (g *gomahawkImpl) NewDBConnectionRequested(t gomahawk.Tomahawk) (gomahawk.DBConnection, error) {
	return nil, nil
}
func (g *gomahawkImpl) NewStreamConnectionRequested(t gomahawk.Tomahawk, uuid string) (gomahawk.StreamConnection, error) {
	return nil, nil
}
func (g *gomahawkImpl) NewStreamConnection(t gomahawk.Tomahawk, sc gomahawk.StreamConnection) error {
	return nil
}
