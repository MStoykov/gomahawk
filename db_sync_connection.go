package gomahawk

import (
	"encoding/json"
	"errors"
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
	m := msg.NewMsg(j, msg.JSON)
	log.Println("sending fetch ops m :[ ", m, "]")
	size, err := d.conn.Write(m.Bytes())

	log.Println("sending fetch ops wrote  ", size, " bytes")
	return err
}

type dBConn struct {
	*secondaryConnection
	offer            *msg.DBsyncOffer
	fom              msg.FetchOpsMethod
	commandProcessor *msg.CommandParser
}

func newDBConn(conn *secondaryConnection) (*dBConn, error) {
	d := new(dBConn)
	d.secondaryConnection = conn
	d.commandProcessor = msg.NewCommandParser()

	go func() {
		for {
			m, err := d.ReadMsg()
			err = d.handleMsg(m)
			if err != nil {
				log.Println(err)
			}
		}
	}()

	return d, nil
}

func openNewDBConn(offer *msg.DBsyncOffer, conn *connection, controlid string) (*dBConn, error) {
	log.Println("OpenNewDBConn")
	d := new(dBConn)
	d.offer = offer
	d.connection = conn
	d.commandProcessor = msg.NewCommandParser()

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
	m := msg.NewMsg(offerBytes, msg.SETUP|msg.JSON)
	log.Println("gonna send msg", m)
	_, err = d.conn.Write(m.Bytes())
	if err != nil {
		log.Println("error while sending offer on dbconnection")
		return nil, err
	}

	go func() {
		for {
			m, err := d.ReadMsg()
			err = d.handleMsg(m)
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
	return d.sendFetchOps(id)
}

func (d *dBConn) handleMsg(m *msg.Msg) error {

	if m.IsDBOP() {
		log.Println("got a DBOP")
		if d.fom != nil {
			command, err := d.commandProcessor.ParseCommand(m)

			if err != nil {
				return err
			}

			err = d.fom.SendCommand(command)

			if !m.IsFragment() {
				return d.fom.Close()
			}

			return err
		} else {
			return errors.New("Got DBOP but no FetchOpsMethod")
		}

	}

	offer, err := msg.ParseDBSyncOffer(m)
	if err != nil {
		log.Println(m)
		return err
	}
	log.Println(offer)

	// parse
	return nil
}
