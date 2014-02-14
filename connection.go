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
			if err != nil {
				c.lastError = err
				return
			}
			err = c.msgHandler(m)
			if err != nil {
				c.lastError = err
				return
			}
		}

	}()
}

func (c *connection) WriteMsg(message *msg.Msg) error {
	_, err := message.WriteTo(c.conn) // actually handle error ?
	return err
}

func (c *connection) ReadMsg() (*msg.Msg, error) {
	return msg.ReadMSG(c.conn)
}

func (c *connection) receiveOffer() error {
	message, err := c.ReadMsg()
	offer, err := msg.ParseOffer(message)
	if err != nil {
		return err
	}
	c.OfferMsg = offer
	if err := c.sendVersionCheck(); err != nil {
		log.Println("versionCheck sending failed for ", offer)

		c.conn.Close()
		return err
	}

	message, err = c.ReadMsg()
	if !message.IsSetup() || !bytes.Equal(message.Payload(), []byte("ok")) {
		log.Println("the version was not ok with the remote for ", offer)
		c.conn.Close()

		return errors.New("Setup failed - wrong version of protocol")
	}

	return nil
}

func (c *connection) sendOffer(offer *msg.Msg) error {
	if err := c.WriteMsg(offer); err != nil {
		c.conn.Close()
		return err
	}

	message, err := c.ReadMsg()
	if err != nil {
		return err
	}

	if !message.IsSetup() || string(message.Payload()) != "4" {
		log.Println("We didn't get versionCheck")
		log.Println("We got", message)
		c.conn.Close()
		return errors.New("VersionCheck failed")
	}

	message = msg.NewMsg([]byte("ok"), msg.SETUP)

	if err = c.WriteMsg(message); err != nil {
		return err
	}

	return nil
}

func (c *connection) Close() error {
	return c.conn.Close()
}

func (c *connection) sendVersionCheck() error {
	message := msg.NewMsg([]byte{'4'}, msg.SETUP)
	return c.WriteMsg(message)
}

type secondaryConnection struct {
	*connection
	parent *controlConnection
}

func newSecondaryConnection(connection *connection, parent *controlConnection) *secondaryConnection {
	return &secondaryConnection{
		connection: connection,
		parent:     parent,
	}
}
