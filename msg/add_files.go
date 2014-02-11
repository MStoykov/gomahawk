package msg

// The AddFiles command says that the files
// listed have been added in this command
type AddFiles struct {
	command
	Files []File `json:"files"`
}
