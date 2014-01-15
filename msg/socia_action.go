package msg

import (
	"encoding/json"
)

// A social Action command
type SocialAction struct {
	Song
	Command
	Action    string `json:"action"`    
	// String representation of the action. "Love" is the only one
	// currently
	Comment   string `json:"comment"` 
	// comment to the action "true" means Love a song "false" Unlove when
	// Action == Love
	Timestamp int64  `json:"timestamp"` // timestamp of the action
}

/*
func (s *SocialAction) GetComment() string {
	return s.Comment
}

func (s *SocialAction) GetTimestamp() int64 {
	return s.Timestamp
}

func (s *SocialAction) GetAction() string {
	return s.Action
}
*/

func NewSocialAction(msg *Msg) (*SocialAction, error) {
	s := new(SocialAction)
	err := json.Unmarshal(msg.Payload(), s)
	if err != nil {
		return nil, err
	}
	return s, nil
}
