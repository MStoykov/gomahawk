package msg

import (
	"encoding/json"
)

// Generic command
type command struct {
	Command string `json:"command"` // "command" : "logplayback",
	Guid    string `json:"guid"`    // "guid" : "33a4a23d-afc4-4b48-9279-b3fd9dff6893",
}

func (c *command) GetCommand() string {
	return c.Command
}
func (c *command) GetGuid() string {
	return c.Guid
}

func WrapCommand(c Command, islast bool) *Msg {
	b, err := json.Marshal(c)
	if err != nil {
		panic(err)
	}

	var flags byte = DBOP | JSON
	if !islast {
		flags |= FRAGMENT
	}
	return NewMsg(b, flags)

}
