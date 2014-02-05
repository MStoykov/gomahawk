package msg

// Rename a playlist
type RenamePlaylist struct {
	command
	PlaylistGuid  string `json:"playlistguid"`
	PlaylistTitle string `json:"playlistTitle"`
}
