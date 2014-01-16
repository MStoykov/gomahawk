package msg

// A playlist
type Playlist struct {
	Info            string `json:"info"`            // generic info
	Creator         string `json:"creator"`         // string represeting the creator
	CreatedOn       int64  `json:"createdon"`       // the time it was created as seconds since 1970-01-01
	Title           string `json:"title"`           // the title of the Playlist
	CurrentRevision string `json:"currentrevision"` // revision number
	Shared          bool   `json:"shared"`          // no idea
	Guid            string `json:"guid"`            // UUID of the Playlist
}
