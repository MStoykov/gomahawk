package logger

import (
	"log"
	"net"
	"sync"
	"time"

	"github.com/MStoykov/gomahawk"
)

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

func ExampleNewGomahawkServer() {
	t := gomahawkImpl{
		"test",
	}
	g, err := gomahawk.NewGomahawkServer(&t)
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
