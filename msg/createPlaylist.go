package msg


// thiscommand says that the given Playlist has been created by this command
type CreatePlaylist struct {
	Command
	Playlist Playlist `json:"playlist"`
}
