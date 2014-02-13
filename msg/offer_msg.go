package msg

import (
	"encoding/json"
	"strconv"
)

// Generic offer Msg payload
type OfferMsg struct {
	Conntype  string `json:"conntype"`            // : "accept-offer"
	Key       string `json:"key"`                 //  "key" : "whitelist",/ "bb3c6870-ac0f-4822-abad-4439e7ffeb15" / "FILE_REQUEST_KEY:12"
	NodeId    string `json:"nodeid,omitempty"`    //  "nodeid" : "bb3c6870-ac0f-4822-abad-4439e7ffeb15",
	ControlId string `json:"controlid,omitempty"` //  "controlid" : "bb3c6870-ac0f-4822-abad-4439e7ffeb15",
	Port      int    `json:"port"`                //  "port" : 0
}

// parse the payload of the given message as an Offer
func ParseOffer(msg *Msg) (*OfferMsg, error) {
	offer := new(OfferMsg)
	err := json.Unmarshal(msg.payload, offer)
	if err != nil {
		offer = nil
	}

	return offer, err
}

// Make new Msg that contains a Request for a File with the given id from the
// given control id
func NewFileRequestOffer(fileId int64, controlid string) *Msg {
	o := &OfferMsg{
		"accept-offer",
		"FILE_REQUEST_KEY:" + strconv.FormatInt(fileId, 10),
		"",
		controlid,
		0,
	}

	r, _ := json.Marshal(o)

	return NewMsg(r, SETUP|JSON)
}

func NewSecondaryOffer(controlid, offerId string, port int) *Msg {
	o := OfferMsg{
		Conntype:  "accept-offer",
		ControlId: controlid,
		Key:       offerId,
		Port:      port,
	}

	r, _ := json.Marshal(o)

	return NewMsg(r, SETUP|JSON)
}
