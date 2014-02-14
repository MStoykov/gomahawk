package gomahawk

import (
	"net"
)

type fakeGomahawk struct{

}

func NewFakeGomahawk() *fakeGomahawk{
	return new(fakeGomahawk)
}
func (f *fakeGomahawk) Name() string {
	return "fake"
}

func (f *fakeGomahawk) ConnectionIsRequested(net.Addr) bool{
	return true
}

func (f *fakeGomahawk) NewTomahawkFound(addr net.Addr, name string) bool {
	return true
}

func (f *fakeGomahawk) NewDBConnection(Tomahawk, DBConnection) (DBConnection, error) {
	return nil, nil
}

func (f *fakeGomahawk) NewDBConnectionRequested(Tomahawk, DBConnection) (DBConnection, error) {
	return nil, nil
}


func (f *fakeGomahawk) NewStreamConnectionRequested(t Tomahawk, uuid string) (StreamConnection, error) {
	return nil, nil
}

