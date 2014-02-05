package msg

// Command saying that the playlist wit the given Guid is deleted
type DeletePlaylist struct {
	command
	PlaylistGuid string `json:"playlistguid"`
}
