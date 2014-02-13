package gomahawk

import (
	"bytes"
	"errors"
	"io"
	"log"

	msg "github.com/MStoykov/gomahawk/msg"
)

type connection struct {
	conn io.ReadWriteCloser
	*msg.OfferMsg
	msgHandler func(*msg.Msg) error
	lastError  error
}

func (c *connection) StartHandelingMessages() {
	go func() { // do it with select and checking whether msgHandler is not nil
		for {
			m, err := c.ReadMsg()
			err = c.msgHandler(m)
			if err != nil {
				c.lastError = err
				return
			}
		}

	}()
}

func (c *connection) ReadMsg() (*msg.Msg, error) {
	return msg.ReadMSG(c.conn)
}

func (c *connection) receiveOffer() error {
	m, err := c.ReadMsg()
	offer, err := msg.ParseOffer(m)
	if err != nil {
		return err
	}
	c.OfferMsg = offer
	if err := c.sendVersionCheck(); err != nil {
		log.Println("versionCheck sending failed for ", offer)

		c.conn.Close()
		return err
	}

	m, err = c.ReadMsg()
	if !m.IsSetup() || !bytes.Equal(m.Payload(), []byte("ok")) {
		log.Println("the version was not ok with the remote for ", offer)
		c.conn.Close()

		return errors.New("Setup failed - wrong version of protocol")
	}

	return nil
}

func (c *connection) sendOffer(offer *msg.Msg) error {
	_, err := offer.WriteTo(c.conn)
	if err != nil {
		c.conn.Close()
		return err
	}

	m, err := c.ReadMsg()
	if !m.IsSetup() || string(m.Payload()) != "4" {
		log.Println("We didn't get versionCheck")
		log.Println("We got", m)
		c.conn.Close()
		return errors.New("VersionCheck failed")
	}

	m = msg.NewMsg([]byte("ok"), msg.SETUP)

	_, err = m.WriteTo(c.conn)
	if err != nil {
		return err
	}

	return nil
}

func (c *connection) Close() error {
	return c.conn.Close()
}

func (c *connection) sendVersionCheck() error {
	m := msg.NewMsg([]byte{'4'}, msg.SETUP)
	_, err := m.WriteTo(c.conn)
	if err != nil {
		return err
	}
	return nil
}

type secondaryConnection struct {
	*connection
	parent *controlConnection
}

func newSecondaryConnection(connection *connection, parent *controlConnection) (*secondaryConnection, error) {
	c := new(secondaryConnection)
	c.connection = connection
	c.parent = parent
	return c, nil
}
