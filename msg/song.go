package msg

// A song 
type Song struct {
	Artist string `json:"artist"` // Name of the artist 
	Track  string `json:"track"`  // Name of the track
}
/*&

func (s *Song) GetArtist() string {
	return s.Artist
}
func (s *Song) GetTrack() string {
	return s.Track
}
*/
