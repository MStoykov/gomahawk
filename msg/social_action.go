package msg

// A social Action command
type SocialAction struct {
	Song
	command
	Action string `json:"action"`
	// String representation of the action. "Love" is the only one
	// currently
	Comment string `json:"comment"`
	// comment to the action "true" means Love a song "false" Unlove when
	// Action == Love
	Timestamp int64 `json:"timestamp"` // timestamp of the action
}
