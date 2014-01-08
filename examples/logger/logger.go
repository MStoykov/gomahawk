package main

import (
	"log"
	"net"
	"sync"
	"time"

	"github.com/MStoykov/gomahawk"
)
type GomahawkImpl struct {
	name string
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
func (g *GomahawkImpl) NewDBConnection(t gomahawk.Tomahawk, db gomahawk.DBConnection) error {
	return nil
}
func (g *GomahawkImpl) NewDBConnectionRequested(t gomahawk.Tomahawk) (gomahawk.DBConnection, error) {
	return nil, nil
}
func (g *GomahawkImpl) NewStreamConnectionRequested(t gomahawk.Tomahawk, uuid string) (gomahawk.StreamConnection, error) {
	return nil, nil
}
func (g *GomahawkImpl) NewStreamConnection(t gomahawk.Tomahawk, sc gomahawk.StreamConnection) error {
	return nil
}
type fetchOpsImpl struct {
}

func (l *fetchOpsImpl) AddFiles(command gomahawk.AddFilesCommand) error {
	log.Println("Got Addfiles with ", len(command.GetFiles()), "files")
	return nil
}

func (l *fetchOpsImpl) DeleteFiles(command gomahawk.DeleteFilesCommand) error {
	log.Println("Got Deletefiles with ", len(command.GetIds()), "ids")
	return nil
}

func (l *fetchOpsImpl) CreatePlaylist(gomahawk.CreatePlaylistCommand) error {
	return nil
}
func (l *fetchOpsImpl) RenamePlaylist(gomahawk.RenamePlaylistCommand) error {
	return nil
}
func (l *fetchOpsImpl) SetPlaylistRevision(gomahawk.SetPlaylistRevisionCommand) error {
	return nil
}

func (l *fetchOpsImpl) DeletePlaylist(gomahawk.DeletePlaylistCommand) error {
	return nil
}

func (l *fetchOpsImpl) SocialAction(gomahawk.SocialActionCommand) error {
	return nil
}

func (l *fetchOpsImpl) Playing(gomahawk.PlayingCommand) error {
	return nil
}

func (l *fetchOpsImpl) StopPlaying(gomahawk.PlayingCommand) error {
	return nil
}

func main() {
	t := new(GomahawkImpl)
	g, err := gomahawk.NewGomahawkServer(t)
	if err != nil {
		log.Println(err)
		return
	}
	err = g.ListenTo(net.IPv4(192, 168, 1, 13), "50210")
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
