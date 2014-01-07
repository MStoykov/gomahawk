package gomahawk

import (
	"log"
	"net"
	"sync"
	"time"
)

type tomahawkImpl struct {
	name string
}

func (t *tomahawkImpl) Name() string{
	return t.name
}

func (t *tomahawkImpl) UUID() string {
	return ""
}

func (t *tomahawkImpl) StreamConnection() (StreamConnection, error) {
	return nil, NotSupportedConnection
}

func (t *tomahawkImpl) DBConnection() (DBConnection, error) {
	return new(loggingDBConnection), nil
}

type loggingDBConnection struct {

}

func (l *loggingDBConnection) AddFiles(command AddFilesCommand) error {
	log.Println("Got Addfiles with ", len(command.GetFiles()), "files")
	return nil
}

func (l *loggingDBConnection) DeleteFiles(command DeleteFilesCommand) error {
	log.Println("Got Deletefiles with ", len(command.GetIds()), "ids")
	return nil
}

func (l *loggingDBConnection) CreatePlaylist(CreatePlaylistCommand) error {
	return nil
}
func (l *loggingDBConnection) RenamePlaylist(RenamePlaylistCommand) error {
	return nil
}
func (l *loggingDBConnection) SetPlaylistRevision(SetPlaylistRevisionCommand) error {
	return nil
}

func (l *loggingDBConnection) DeletePlaylist(DeletePlaylistCommand) error {
	return nil
}

func (l *loggingDBConnection) Love(LoveCommand) error {
	return nil
}

func (l *loggingDBConnection) UnLove(LoveCommand) error  {
	return nil
}

func (l *loggingDBConnection) Playing(PlayingCommand) error{
	return nil
}

func (l *loggingDBConnection) StopPlaying(PlayingCommand) error{
	return nil
}

func (l *loggingDBConnection) FetchOps(string) error{
	return nil
}
func ExampleNewGomahawk() {
	t := tomahawkImpl{
		"test",
	}
	g, err := NewGomahawk(&t)
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
