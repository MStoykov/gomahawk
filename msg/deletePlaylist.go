package msg

// Command saying that the playlist wit the given Guid is deleted
type DeletePlaylist struct {
	Command
	PlaylistGuid string `json:"playlistguid"`
}
