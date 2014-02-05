package msg

// Change the playlist
//
// Only new entries are being given in AddEntries
// the OrderedGuids is the order guids of all the song in the new revision of
// the playlist. To remove a song you simply skip it in the OrderedGuids
type SetPlaylistRevision struct {
	command
	OldRev       string     `json:"oldrev"`       // the old revision of the playlist
	NewRev       string     `json:"newrev"`       // the new revision of the playlist
	PlaylistGuid string     `json:"playlistguid"` // the Guid of the Playlist
	AddEntries   []AddEntry `json:"addedentries"` // the entries to be added
	OrderedGuids []string   `json:"orderedguids"` // the ordered list Of All the songs in the playlist
}

type AddEntry struct {
	Duration     int64     `json:"duration"`     // the duration of the song
	LastModified int64     `json:"lastmodified"` //
	Guid         string    `json:"guid"`         // the guid of the song
	Annotation   string    `json:"annotation"`   // a annotation ?
	Query        QueryType `json:"query"`        // still not sure
}

type QueryType struct {
	Song
	Album    string `json:"album"`    // the album of the song
	Duration int    `json:"duration"` // duration of the song (though it's -1 ) allthe time
	Qid      string `json:"qid"`      // some ID
}
