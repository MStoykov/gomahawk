package gomahawk

import (
	"bytes"
	"errors"
	"io"
	"log"

	msg "github.com/MStoykov/gomahawk/msg"
)

func (d *dBConn) Trigger() error {
	return nil
}

// request changes after given id. "" means all
func (d *dBConn) sendFetchOps(id string) error {
	m := msg.NewFetchOpsMsg(id)
	_, err := m.WriteTo(d.conn)
	return err
}

type dBConn struct {
	*secondaryConnection
	offer            *msg.DBsyncOffer
	dbconn           DBConnection
	fom              msg.FetchOpsMethod
	commandProcessor *msg.CommandParser
}

func newDBConn(conn *secondaryConnection) (*dBConn, error) {
	d := new(dBConn)
	d.secondaryConnection = conn
	d.commandProcessor = msg.NewCommandParser()

	d.msgHandler = d.handleMsg

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
	_, err := m.WriteTo(d.conn)
	if err != nil {
		log.Println("error while sending offer on dbconnection")
		return nil, err
	}

	d.msgHandler = d.handleMsg

	return d, nil
}

func (d *dBConn) FetchOps(fom msg.FetchOpsMethod, id string) error {
	d.fom = fom
	return d.sendFetchOps(id)
}

func (d *dBConn) handleMsg(m *msg.Msg) error {
	if m.IsDBOP() {
		if d.fom == nil {
			return errors.New("Got DBOP but no FetchOpsMethod")
		}

		if !m.IsJSON() { // 'ok' ?
			if bytes.Equal(m.Payload(), []byte("ok")) {
				return d.fom.Close()
			} else {
				return errors.New("Got DBOP that's JSON but not 'ok' :\n" + m.String())
			}
		}

		command, err := d.commandProcessor.ParseCommand(m)

		if err != nil {
			if nerr, ok := err.(msg.NotRegisteredError); ok {
				log.Println(nerr)
				return nil
			}
			return err
		}

		err = d.fom.SendCommand(command)

		if !m.IsFragment() {
			return d.fom.Close()
		}

		return err
	}

	op, err := msg.GetOpFromFetchOpsMsg(m)

	if err == nil {
		if d.dbconn != nil {
			return d.dbconn.FetchOps(newDummyFetchOps(d.conn), op)
		} else {
			_, err := msg.NewMsg([]byte("ok"), msg.DBOP).WriteTo(d.conn)
			return err
		}
	}
	if msg.IsTrigger(m) {
		if d.dbconn != nil {
			return d.dbconn.Trigger()
		}

		return nil
	}

	return nil
}

type DummyFetchOps struct {
	lastCommand msg.Command
	writer      io.Writer
}

func newDummyFetchOps(writer io.Writer) msg.FetchOpsMethod {
	return &DummyFetchOps{
		writer: writer,
	}
}

func (d *DummyFetchOps) SendCommand(command msg.Command) (err error) {
	if d.lastCommand != nil {
		_, err = msg.WrapCommand(d.lastCommand, false).WriteTo(d.writer)
	}

	d.lastCommand = command

	return
}

func (d *DummyFetchOps) Close() (err error) {
	if d.lastCommand != nil {
		_, err = msg.WrapCommand(d.lastCommand, false).WriteTo(d.writer)
	} else {
		_, err = msg.NewMsg([]byte("ok"), msg.SETUP).WriteTo(d.writer)
	}

	return
}
