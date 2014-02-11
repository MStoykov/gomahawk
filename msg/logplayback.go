package msg

// A command saying that a given song has been started or stopped playing
// action == 1 is start playing, action == 2 is stopped playing
// the GUID of an Logplayback with action == 1 can not be used in fetchops
type LogPlayback struct {
	Song
	command
	Action     int   `json:"action"`     // 1 for start ,2  for stop
	PlayTime   int64 `json:"playtime"`   // time in seconds since 1970-01-01
	SecsPlayed int   `json:"secsPlayed"` // seconds of the song that have been played
}
