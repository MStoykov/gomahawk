package gomahawk

import (
	"bytes"
	"errors"
	"log"

	msg "github.com/MStoykov/gomahawk/msg"
)

func (d *dBConn) Trigger() error {
	return nil
}

// request changes after given id. "" means all
func (d *dBConn) sendFetchOps(id string) error {
	return d.WriteMsg(msg.NewFetchOpsMsg(id))
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

	message := msg.NewSecondaryOffer(controlid, offer.Key, 50210)
	if err := d.WriteMsg(message); err != nil {
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

func (d *dBConn) handleMsg(message *msg.Msg) error {
	if message.IsDBOP() {
		if d.fom == nil {
			return errors.New("Got DBOP but no FetchOpsMethod")
		}

		if !message.IsJSON() { // 'ok' ?
			if bytes.Equal(message.Payload(), []byte("ok")) {
				return d.fom.Close()
			} else {
				return errors.New("Got DBOP that's JSON but not 'ok' :\n" + message.String())
			}
		}

		command, err := d.commandProcessor.ParseCommand(message)

		if err != nil {
			if nerr, ok := err.(msg.NotRegisteredError); ok {
				log.Println(nerr)
				return nil
			}
			return err
		}

		err = d.fom.SendCommand(command)

		if !message.IsFragment() {
			return d.fom.Close()
		}

		return err
	}

	op, err := msg.GetOpFromFetchOpsMsg(message)

	if err == nil {
		if d.dbconn != nil {
			return d.dbconn.FetchOps(newDummyFetchOps(d.connection), op)
		} else {
			return d.WriteMsg(msg.NewMsg([]byte("ok"), msg.DBOP))
		}
	}

	if msg.IsTrigger(message) {
		if d.dbconn != nil {
			return d.dbconn.Trigger()
		}

		return nil
	}
	log.Println("unhandled message received on a DBConnection\n", message)

	return nil
}

type DummyFetchOps struct {
	lastCommand msg.Command
	conn        *connection
}

func newDummyFetchOps(conn *connection) msg.FetchOpsMethod {
	return &DummyFetchOps{
		conn: conn,
	}
}

func (d *DummyFetchOps) SendCommand(command msg.Command) (err error) {
	if d.lastCommand != nil {
		err = d.conn.WriteMsg(msg.WrapCommand(d.lastCommand, false))
	}

	d.lastCommand = command

	return
}

func (d *DummyFetchOps) Close() (err error) {
	if d.lastCommand != nil {
		err = d.conn.WriteMsg(msg.WrapCommand(d.lastCommand, false))
	} else {
		err = d.conn.WriteMsg(msg.NewMsg([]byte("ok"), msg.SETUP))
	}

	return
}
