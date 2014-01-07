// Implementation of the network protocol of tomahawk
//
// Basic Usage :
// 
// 1.Implement Tomahawk
//
// 2.Make newInstance of Gomahawk with NewGomahawk(TomahawkImpl)
//
// 3.Gomahawk.Start()
package gomahawk

import (
		"errors"
)

var (

	NotSupportedConnection = errors.New("Not Supported Connection")
)

// This is the Interface that represents a remote Tomahawk as well as the local. 
//
// As user of the library this is the interface that needs to be implemented.
type Tomahawk interface {
	// Returns a human readable name of the instance. Can be empty string
	Name() string
	// Returns the uuid of the tomahawk instance as a string.
	// 
	// While implementing this interface you can use "" and Gomahawk will create a uuid based on the name
	UUID() string
	// Get the streamConnection associated with this tomahawk. 
	//
	// NotSupportedConnection is returned when this type of connection is not supported
	// Others error mean that the connection was not successfully created
	StreamConnection() (StreamConnection, error)
	// Get the DBConnection associated with this tomahawk
	//
	// NotSupportedConnection is returned when this type of connection is not supported
	// Others error mean that the connection was not successfully created
	DBConnection() (DBConnection, error)
}

// The StreamConnection is used to transfer files over the network
type StreamConnection interface {
	// returns the size of each block
	// this should be a constant for each instance
	BlockSize() int
	// Returns the blockIndex-ed block of file with the given id 
	// 
	// For all but the last block it's required that the result is of length BlockSize
	// And for any blockIndex > 0 this function should not panic or return a nil slice
	Block(blockIndex int, id int64) []byte
}

// DBConnection are used to show the database to another tomahawk instance and to sync it afterwards
// 
// The order in which the functions are called is the order in which the commands are ordered
type DBConnection interface {
	AddFiles(AddFilesCommand) error
	DeleteFiles(DeleteFilesCommand) error
	CreatePlaylist(CreatePlaylistCommand) error
	RenamePlaylist(RenamePlaylistCommand) error
	SetPlaylistRevision(SetPlaylistRevisionCommand) error
	DeletePlaylist(DeletePlaylistCommand) error
	Love(LoveCommand) error
	UnLove(LoveCommand) error 
	Playing(PlayingCommand) error
	StopPlaying(PlayingCommand) error
	// request that all changes since the given UUID of a command are send.
	// "" means all changes
	FetchOps(string) error
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

type LoveCommand interface {
	Command
	Song
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
