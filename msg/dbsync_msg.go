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
		offer = nil
	}
	if offer.Method != "dbsync-offer" {
		return nil, errors.New("not DBSYNCOFFER")
	}

	return offer, err
}
