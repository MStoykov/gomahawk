package msg

import (
	"encoding/json"
	"errors"
)

// the message send forgetting dbsync
type DBsyncOffer struct {
	Method string `json:"method"` // "dbsync-offer"
	Key    string `json:"key"`    // uuid
}

func ParseDBSyncOffer(msg *Msg) (*DBsyncOffer, error) {
	offer := new(DBsyncOffer)
	err := json.Unmarshal(msg.payload, offer)
	if err != nil {
		return nil, err
	}
	if offer.Method != "dbsync-offer" {
		return nil, errors.New("not DBSYNCOFFER")
	}

	return offer, nil
}

func NewDBSyncOfferMsg(key string) (m *Msg) {
	offer := DBsyncOffer{
		Method: "dbsync-offer",
		Key:    key,
	}
	offerBytes, _ := json.Marshal(offer)

	return NewMsg(offerBytes, SETUP|JSON)
}

func IsTrigger(msg *Msg) bool {
	if msg.IsJSON() {
		msg.Uncompress()
		var m map[string]string
		if err := json.Unmarshal(msg.payload, &m); err != nil {
			return false
		}

		return len(m) == 1 && m["method"] == "trigger"
	}
	return false
}
