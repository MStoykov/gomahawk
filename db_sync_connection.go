package gomahawk

import (
	"bytes"
	"encoding/json"
	"log"
	"regexp"

	msg "github.com/MStoykov/gomahawk/msg"
)

type secondaryOffer struct {
	ConnType  string `json:"conntype"`  // "conntype": "accept-offer"
	ControlId string `json:"controlid"` //"controlid": "66bd135d-113f-481a-977e-111111111111",
	Key       string `json:"key"`       //"key": KEY,
	Port      int    `json:"port"`      // "port": 50210
}

type fetchOpsMethodMsg struct {
	Method string `json:"method"` // method: fetchOps
	Lastop string `json:"lastop"` // lastop :"66bd135d-113f-481a-977e-111111111111"
}

var re = regexp.MustCompile(`"command"\s*:\s*"([^"]+)"`)

func (d *dBConn) Trigger() error {
	return nil
}

func parseFetchOpsMethod(m *msg.Msg) (*fetchOpsMethodMsg, error) {
	r := new(fetchOpsMethodMsg)
	err := json.Unmarshal(m.Payload(), r)
	if err != nil {
		r = nil
	}
	return r, err

}

// request changes after given id. "" means all
func (d *dBConn) sendFetchOps(id string) error {
	met := &fetchOpsMethodMsg{
		Method: "fetchops",
		Lastop: id,
	}
	j, err := json.Marshal(met)
	if err != nil {
		return err
	}
	log.Println(j)
	m := msg.NewMsg(bytes.NewBuffer(j), msg.JSON)
	log.Println("sending fetch ops m :[ ", m, "]")
	size, err := d.conn.Write(m.Bytes())

	log.Println("sending fetch ops wrote  ", size, " bytes")
	return err
}

type dBConn struct {
	*secondaryConnection
	offer *msg.DBsyncOffer
	fom   msg.FetchOpsMethod
}

func newDBConn(conn *secondaryConnection) (*dBConn, error) {
	d := new(dBConn)
	d.secondaryConnection = conn

	//	go func() {
	//		for m := range d.sync {
	//			err := d.handleMsg(m)
	//			if err != nil {
	//				log.Println(err)
	//				return
	//			}
	//		}
	//	}()

	return d, nil
}

func openNewDBConn(offer *msg.DBsyncOffer, conn *connection, controlid string) (*dBConn, error) {
	log.Println("OpenNewDBConn")
	d := new(dBConn)
	d.offer = offer
	d.connection = conn

	dbsecondaryoffer := secondaryOffer{
		"accept-offer",
		controlid,
		offer.Key,
		50210, //hardcoded
	}
	offerBytes, err := json.Marshal(dbsecondaryoffer)

	if err != nil {
		log.Println("error while marshaling offer", err)
		return nil, err
	}
	m := msg.NewMsg(bytes.NewBuffer(offerBytes), msg.SETUP|msg.JSON)
	log.Println("gonna send msg", m)
	_, err = d.conn.Write(m.Bytes())
	if err != nil {
		log.Println("error while sending offer on dbconnection")
		return nil, err
	}

	sync := make(chan *msg.Msg)
	d.sync = sync
	d.processor = msg.NewProcessor(d.conn, sync)

	go func() {
		for m := range d.sync {
			err := d.handleMsg(m)
			if err != nil {
				log.Println(err)
				return
			}
		}
	}()

	return d, nil
}

func (d *dBConn) FetchOps(fom msg.FetchOpsMethod, id string) error {
	d.fom = fom
	d.sendFetchOps(id)
	go func() {
		others := make(chan *msg.Msg)
		go func() {
			for m := range others {
				d.handleMsg(m)
			}
		}()
		err := msg.FilterCommands(d.sync, others, d.fom)
		if err != nil {
			log.Println("error from FilterCommands", err)
		} else {
			log.Println("no error on FilterCommands")
		}

		d.fom = nil
	}()
	return nil
}

func (d *dBConn) handleMsg(m *msg.Msg) error {

	offer, err := msg.ParseDBSyncOffer(m)
	if err != nil {
		return err
	}
	log.Println(offer)

	// parse
	return nil
}
