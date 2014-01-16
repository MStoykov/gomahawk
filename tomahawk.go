package gomahawk

import (
	"time"

	. "github.com/MStoykov/gomahawk/msg"
)

// This is the Interface that represents a remote Tomahawk.
type Tomahawk interface {
	// Returns a human readable name of the instance. Can be empty string
	Name() string
	// Returns the uuid of the tomahawk instance as a string.
	//
	// While implementing this interface you can use "" and a uuid based on the name
	UUID() string
	// Get a streamConnection for the file with the given id previously received from AddFiles command from this tomahawk.
	//
	// NotSupportedConnection is returned when this type of connection is not supported
	// Others error mean that the connection was not successfully created
	RequestStreamConnection(uuid int64) (StreamConnection, error)
	// A non blocking request for a DBConnection
	// a call to NewDBConnection with the associated Tomahawk and the DBConnection will be made
	// when and if the remote tomahawks agrees to it.
	// Multiple calls to this
	//
	// NotSupportedConnection is returned when this type of connection is not supported
	// Others error mean that the connection was not successfully created
	RequestDBConnection(DBConnection) error

	// The last time a ping message was received from this particular Tomahawk
	LastPing() time.Time
}

// The StreamConnection is used to transfer files over the network
type StreamConnection interface {
	// returns the ID of the file associated with this connection
	FileID() int64
	// return the current stream. Multiple calls to this function without SeekToBlock
	// return the same channel
	GetStream() (<-chan []byte, error)
	// this will close the previos stream and make a new one that will start giving blocks from the given index forward
	// the size of block is not set but it should be considered constant for each StreamConnection
	//
	// The original Tomahawk has it hardcoded to 4096 bytes at the time of writing
	SeekToBlock(blockIndex int) error
}

type DBConnection interface {
	// Signals that there are changes to the Database since the last FetchOps
	// Should not be called while on the
	Trigger() error
	// request that all changes since the given UUID of a command are send using the given FetchOpsMethod.
	//
	// You can not get or make more then one simultanious FetchOps.
	//
	// "" means all changes otherwise it is a uuid of previously transfered command. Look at the protocol specification for which command's uuids can be used.
	FetchOps(FetchOpsMethod, string) error

	// Close the underlying connection.
	Close() error
}
