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
	"time"
)

var (
	NotSupportedConnection = errors.New("Not Supported Connection")
)

// This is the Interface that represents a remote Tomahawk.
type Tomahawk interface {
	// Returns a human readable name of the instance. Can be empty string
	Name() string
	// Returns the uuid of the tomahawk instance as a string.
	//
	// While implementing this interface you can use "" and a uuid based on the name
	UUID() string
	// Get a streamConnection for the file with the given uuid with this tomahawk.
	//
	// NotSupportedConnection is returned when this type of connection is not supported
	// Others error mean that the connection was not successfully created
	RequestStreamConnection(uuid string) error
	// A non blocking request for a DBConnection
	// a call to NewDBConnection with the associated Tomahawk and the DBConnection will be made
	// when and if the DBConnection is made
	//
	// NotSupportedConnection is returned when this type of connection is not supported
	// Others error mean that the connection was not successfully created
	RequestDBConnection() error

	// Signals to the Tomahawk that there are changes to the db
	// This may result in a db connection being requested
	TriggerDBChanges()

	// The last time a ping message was received from this particular Tomahawk
	LastPing() time.Time
}

// The StreamConnection is used to transfer files over the network
type StreamConnection interface {
	// returns the ID of the file associated with this connection
	FileID() int64
	// returns the size of each block in bytes
	//
	// this should be a constant for each instance
	BlockSize() int
	// returns a channel of []byte each of which represent a Block of the file being requested.
	// The blockIndex is from where the streaming should begin.
	//
	// When the channel is closed it means that the last block has been sent.
	//
	// Consequentive calls to this function invalidate previous channel, which will be closed.
	StartFromBlock(blockIndex int) (<-chan []byte, error)
}

type DBConnection interface {
	// request that all changes since the given UUID of a command are send using the given FetchOpsMethod.
	//
	// You can not get or make more then one simultanious FetchOps.
	//
	// "" means all changes otherwise it is a uuid of previously transfered command
	FetchOps(FetchOpsMethod, string) error
}

// This is interface represents ONE fetchop workflow
//
// Every call means a given command has been emmited. The commands are in order and the interface CAN NOT be used in parallel.
//
// Returning from a method means that the next command can be received. Multiple simultanius calls to the inteface may cause panic.
//
// None of the commands are guaranteed to be send before Close is called because of the way the protocol is designed.
// And every last call is never send before Close is called. After Close or after any error returned from any function it should be considered that ALL
// commands were not send.
type FetchOpsMethod interface {
	AddFiles(AddFilesCommand) error
	DeleteFiles(DeleteFilesCommand) error
	CreatePlaylist(CreatePlaylistCommand) error
	RenamePlaylist(RenamePlaylistCommand) error
	SetPlaylistRevision(SetPlaylistRevisionCommand) error
	DeletePlaylist(DeletePlaylistCommand) error
	SocialAction(SocialActionCommand) error
	Playing(PlayingCommand) error
	StopPlaying(PlayingCommand) error
	Close() error
}

type Command interface {
	// uuid that is unique for each command
	Guid() string
}

type Song interface {
	// the artist of the song
	Artist() string
	// the name of the song
	Track() string
}

type PlayingCommand interface {
	Command
	Song
	// how long the track is in seconds
	TrackDuration() int
	// how many seconds were played
	PlayedSeconds() int
	// the time as seconds since 1970-01-01
	Playtime() int64
}

// SocialAction is currently used for signaling love and unloving of a song.
// if Action is equal to "Love" then Comment is should be either "true" or "false" -
// respectfully this signals loving the giving song or unloving it
type SocialActionCommand interface {
	Command
	Song
	Action() string
	Comment() string
	Timestamp() int64
}

type AddFilesCommand interface {
	Command
	// returns a slice of files that have been added
	GetFiles() []File
}

type DeleteFilesCommand interface {
	Command
	// slice with the ids of the files that have been deleted
	GetIds() []int64
}

type CreatePlaylistCommand interface {
	Command
	Playlist() Playlist
}
type SetPlaylistRevisionCommand interface {
	Command
	// TODO
}

type Playlist interface {
	// General info
	Info() string
	// String represeing the creator of the playlist
	Creator() string
	// The time at which the playlist was created as seconds since 1970-01-01
	CreatedOn() int64
	// The title of the Playlist
	Title() string
	// Revision of this particular instance of the playlist
	CurrentRevision() string
	// Whether the playlist is shared
	Shared() bool
	// the UUID of the playlist as string
	Guid() string
}

type DeletePlaylistCommand interface {
	Command
	// the uuid of the playlist to be deleted
	PlaylistGuid() string
}
type RenamePlaylistCommand interface {
	Command
	// the uuid of the playlist to be renamed
	PlaylistGuid() string
	// the new name of the playlist
	PlaylistTitle() string
}

// A self-explenatory file interface
type File interface {
	// Unique identifier for a track
	ID() int64
	// the same as ID
	URL() string
	// The Artist of the song
	Artist() string
	// The Album the song is from
	Album() string
	// The name of the Song
	Track() string
	// The mimetype of the file
	MimeType() string
	// empty - reserved for future usages
	Hash() string
	// The year the Song/Album was released
	Year() int
	// The position of the song in the Album
	AlbumPos() int
	// The last time the file was changed in seconds since 1970-01-01
	MTime() int64
	// The length of the song in seconds
	Duration() int
	// The bitrate of the song (can be 0)
	Bitrate() int
	// The size of the song in kbytes
	Size() int64
}
