package msg

import(
	"encoding/json"
)
type fetchOpsMethod struct {
	Method string `json:"method"` // method: fetchOps
	LastOp string `json:"lastop"` // lastop :"66bd135d-113f-481a-977e-111111111111"
}

func NewFetchOpsMsg(op string) *Msg{
	f := fetchOpsMethod{
		Method:"fetchops",
		LastOp: op,
	}
	b, _ :=	json.Marshal(f)

	return NewMsg(b, JSON)
}

func GetOpFromFetchOpsMsg(m *Msg) (string, error) {
	f := new(fetchOpsMethod)
	err := json.Unmarshal(m.payload, f)
	if err != nil {
		return "", err
	}

	return f.LastOp, nil
}
