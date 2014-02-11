package msg

type SetCollectionAttributes struct {
	command
	Del  bool   `json:"del"`
	Id   string `json:"id"`
	Type int    `json:"type"`
}
