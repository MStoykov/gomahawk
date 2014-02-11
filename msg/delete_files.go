package msg

// Command saying the the listed ids are of files that have been deleted
type DeleteFiles struct {
	command
	Ids []int64 `json:"ids"`
}
