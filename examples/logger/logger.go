package main

import (
	"log"
	"net"
	"os"
	"sync"

	"github.com/MStoykov/gomahawk"
	msg "github.com/MStoykov/gomahawk/msg"
)

type GomahawkImpl struct {
	name   string
	server gomahawk.GomahawkServer
}

type DBConnectionImpl struct {
	g          *GomahawkImpl
	t          gomahawk.Tomahawk
	other      gomahawk.DBConnection
	outgoingFO msg.FetchOpsMethod
	incomingFO msg.FetchOpsMethod
	lastId     string
}

func (d *DBConnectionImpl) Trigger() error {
	log.Println("Trigger received")
	if d.incomingFO == nil {
		return d.GetSinceLast()
	}

	return nil
}

func (d *DBConnectionImpl) GetSinceLast() error {
	d.incomingFO = &fetchOpsImpl{d}

	return d.other.FetchOps(d.incomingFO, d.lastId)
}

func (d *DBConnectionImpl) Close() error {
	return nil // NOOP
}

func (d *DBConnectionImpl) FetchOps(outgoing msg.FetchOpsMethod, op string) error {
	log.Println("FetchOps since ", op, "received")
	d.outgoingFO = outgoing
	go func() {
		log.Println("we don't have anything so we are gonna close")
		d.outgoingFO.Close()
	}()

	return nil
}

func (g *GomahawkImpl) Name() string {
	return g.name
}
func (g *GomahawkImpl) ConnectionIsRequested(addr net.Addr) bool {
	return true
}
func (g *GomahawkImpl) NewTomahawkFound(addr net.Addr, name string) bool {
	return true
}
func (g *GomahawkImpl) NewDBConnection(t gomahawk.Tomahawk, db gomahawk.DBConnection) (gomahawk.DBConnection, error) {
	log.Printf("NewDBConneciton(%#v, %#v) ", t, db)
	d := &DBConnectionImpl{g: g, t: t, other: db}
	go d.GetSinceLast()
	return d, nil
}
func (g *GomahawkImpl) NewDBConnectionRequested(t gomahawk.Tomahawk, db gomahawk.DBConnection) (gomahawk.DBConnection, error) {
	return nil, gomahawk.NotSupportedConnection
}
func (g *GomahawkImpl) NewStreamConnectionRequested(t gomahawk.Tomahawk, uuid string) (gomahawk.StreamConnection, error) {
	return nil, gomahawk.NotSupportedConnection
}
func (g *GomahawkImpl) NewStreamConnection(t gomahawk.Tomahawk, sc gomahawk.StreamConnection) error {
	return nil
}

type fetchOpsImpl struct {
	dbConn *DBConnectionImpl
}

func (l *fetchOpsImpl) SendCommand(command msg.Command) error {
	log.Println("command received : ", command)
	if command.GetCommand() == "logplayback" {
		if command.(*msg.LogPlayback).Action == 1 { // start playback is not fetchopable
			return nil
		}
	}
	l.dbConn.lastId = command.GetGuid()
	return nil
}

func (l *fetchOpsImpl) Close() error {
	log.Println("fetchops closed. Last command was with uuid", l.dbConn.lastId)
	l.dbConn.incomingFO = nil

	return nil
}

func main() {
	if len(os.Args) != 2 {
		log.Fatal("usage : go logger.go <ip to listen to>")
	}
	ipstring := os.Args[1]
	t := new(GomahawkImpl)
	t.name = "awesome"
	g, err := gomahawk.NewGomahawkServer(t)
	if err != nil {
		log.Println(err)
		return
	}
	t.server = g
	err = g.ListenTo(net.ParseIP(ipstring), 50210)
	if err != nil {
		log.Println(err)
		return
	}

	g.AdvertiseEvery(5)

	err = g.Start()
	if err != nil {
		log.Println(err)
		return
	}
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}
