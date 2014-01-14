package main

import (
	"log"
	"net"
	"sync"
	"time"
	"os"

	"github.com/MStoykov/gomahawk"
	msg "github.com/MStoykov/gomahawk/msg"
)

type GomahawkImpl struct {
	name string
	server 	gomahawk.GomahawkServer
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
	db.FetchOps(&fetchOpsImpl{g,t, "", 1}, "")
	return nil,nil
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
	g *GomahawkImpl
	t gomahawk.Tomahawk
	last string
	fileId int64
}

func (l *fetchOpsImpl) AddFiles(command *msg.AddFiles) error {
	log.Println("Got Addfiles with ", len(command.Files), "files")
	l.fileId = command.Files[0].Id
	l.last = command.Guid
	return nil
}

func (l *fetchOpsImpl) DeleteFiles(command *msg.DeleteFiles) error {
	log.Println("Got Deletefiles with ", len(command.Ids), "ids")
	l.last = command.Guid
	return nil
}

func (l *fetchOpsImpl) CreatePlaylist(command *msg.CreatePlaylist) error {
	l.last = command.Guid
	return nil
}
func (l *fetchOpsImpl) RenamePlaylist(command *msg.RenamePlaylist) error {
	l.last = command.Guid
	return nil
}
func (l *fetchOpsImpl) SetPlaylistRevision(command *msg.SetPlaylistRevision) error {
	l.last = command.Guid
	return nil
}

func (l *fetchOpsImpl) DeletePlaylist(command *msg.DeletePlaylist) error {
	l.last = command.Guid
	return nil
}

func (l *fetchOpsImpl) SocialAction(command *msg.SocialAction) error {
	l.last = command.Guid
	return nil
}

func (l *fetchOpsImpl) LogPlayback(command *msg.LogPlayback) error {
//	log.Println("got stopplaying")
	log.Printf("%#v\n", command)
	l.last = command.Guid
	return nil
}

func (l *fetchOpsImpl) Close() error {
	log.Println("fetchops closed. Last command was with uuid", l.last)
/*log.Println("gonna try to stream a file NOW with id ", l.fileId)

	sc, err :=  l.t.RequestStreamConnection(l.fileId)
	if err != nil {
		log.Println(err)
		return nil
	}
	log.Println("StreamConnection", sc)
	log.Println("sc.FileID()" , sc.FileID())
	c, err := sc.GetStream()
	if err != nil {
		log.Println(err)
	}

	for b := range c {
		log.Println("we read ", len(b), "bytes from the stream")
	}
	
*/
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
	err = g.ListenTo(net.ParseIP(ipstring), "50210")
	if err != nil {
		log.Println(err)
		return
	}

	g.AdvertEvery(time.Second * 10)

	err = g.Start()
	if err != nil {
		log.Println(err)
		return
	}
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}
