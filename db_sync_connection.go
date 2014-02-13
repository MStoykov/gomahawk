package gomahawk

import (
	"errors"
	"log"

	msg "github.com/MStoykov/gomahawk/msg"
)

func (d *dBConn) Trigger() error {
	return nil
}
// request changes after given id. "" means all
func (d *dBConn) sendFetchOps(id string) error {
	m := msg.NewFetchOpsMsg(id)
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

	m := msg.NewSecondaryOffer(controlid, offer.Key, 50210)
	log.Println("gonna send msg", m)
	_, err := d.conn.Write(m.Bytes())
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
