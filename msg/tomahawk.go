package msg


// This is interface represents ONE fetchop workflow
//
// Every call means a given command has been sent. The commands are in order and the interface CAN NOT be used in parallel.
//
// Returning from a method means that the next command can be received. Multiple simultaneous calls to the inteface may cause panic.
//
// None of the commands are guaranteed to be send before Close is called because of the way the protocol is designed.
// And every last call is never send before Close is called. After Close or after any error returned from any function it should be considered that ALL
// commands were not send.
// 
// When the interface is being implemented in order to receive commands returning error will stop the calling of methods but the network activity may still continue
type FetchOpsMethod interface {
	AddFiles(*AddFiles) error
	DeleteFiles(*DeleteFiles) error
	CreatePlaylist(*CreatePlaylist) error
	RenamePlaylist(*RenamePlaylist) error
	SetPlaylistRevision(*SetPlaylistRevision) error
	DeletePlaylist(*DeletePlaylist) error
	SocialAction(*SocialAction) error
	LogPlayback(*LogPlayback) error
	Close() error
}
