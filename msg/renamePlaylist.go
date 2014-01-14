package msg

// Rename a playlist
type RenamePlaylist struct {
	Command
	PlaylistGuid  string `json:"playlistguid"` 
	PlaylistTitle string `json:"playlistTitle"`
}
