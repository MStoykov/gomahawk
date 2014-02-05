package msg

import (
	"encoding/json"
)

// Command saying the the listed ids are of files that have been deleted
type DeleteFiles struct {
	command
	Ids []int64 `json:"ids"`
}

func NewDeleteFiles(m *Msg) (d *DeleteFiles, err error) {
	d = new(DeleteFiles)
	err = json.Unmarshal(m.Payload(), d)
	if err != nil {
		d = nil
	}

	return
}
