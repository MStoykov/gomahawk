package msg

// A File
type File struct {
	Id  int64  `json:"id"`  // the ID of the file
	Url string `json:"url"` // currently only a string representation of the ID
	Song
	Album    string `json:"album"`    // string representation of the Album name
	Mimetype string `json:"mimetype"` // the mime type of the file
	Hash     string `json:"hash"`     // future expanstion
	Year     int    `json:"year"`     // year the song has been release
	Albumpos int    `json:"albumpos"` // position in the album of the song
	Mtime    int64  `json:"mtime"`    // the last modification time of the file in seconds since 1970-01-01
	Duration int    `json:"duration"` // duration of the song in seconds
	Bitrate  int    `json:"bitrate"`  // bitrate of the file (can be 0)
	Size     int    `json:"Size"`     // size of the file in bytes
}

/*
func (f *File) GetId() int64 {
	return f.Id
}

func (f *File) GetUrl() string {
	return f.Url
}

func (f *File) GetAlbum() string {
	return f.Album
}

func (f *File) GetMimetype() string {
	return f.Mimetype
}

func (f *File) GetHash() string {
	return f.Hash
}

func (f *File) GetYear() int {
	return f.Year
}

func (f *File) GetAlbumpos() int {
	return f.Albumpos
}

func (f *File) GetMtime() int64 {
	return f.Mtime
}

func (f *File) GetDuration() int {
	return f.Duration
}

func (f *File) GetBitrate() int {
	return f.Bitrate
}

func (f *File) GetSize() int {
	return f.Size
}
*/
