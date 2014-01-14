package msg

import (
	"encoding/json"
)

// The AddFiles command says that the files 
// listed have been added in this command
type AddFiles struct {
	Command
	Files []File `json:"files"`
}

/*
func (a *AddFiles) GetFiles() []File {
	return a.Files
}
*/

func NewAddFiles(msg *Msg) (*AddFiles, error) {
	a := new(AddFiles)
	err := json.Unmarshal(msg.Payload(), a)
	if err != nil {
		return nil, err
	}
	return a, nil
}
