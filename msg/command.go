package msg


// Generic command
type Command struct {
	Command string `json:"command"` // "command" : "logplayback",
	Guid    string `json:"guid"`    // "guid" : "33a4a23d-afc4-4b48-9279-b3fd9dff6893",
}

/*
func (c *Command) GetCommand() string {
	return c.Command
}
func (c *Command) GetGuid() string {
	return c.Guid
}
*/
