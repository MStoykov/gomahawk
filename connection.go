package gomahawk

import (
	"bytes"
	"errors"
	"io"
	"log"

	msg "github.com/MStoykov/gomahawk/msg"
)

type connection struct {
	conn      io.ReadWriteCloser
	processor *msg.Processor
	sync      <-chan *msg.Msg
	*msg.OfferMsg
}

func (c *connection) setupProcessor() {
	sync := make(chan *msg.Msg)
	c.sync = sync
	c.processor = msg.NewProcessor(c.conn, sync)
}

func (c *connection) receiveOffer() error {
	m := <-c.sync
	offer, err := msg.ParseOffer(m)
	if err != nil {
		return err
	}
	c.OfferMsg = offer
	if err := c.sendversionCheck(); err != nil {
		log.Println("versionCheck sending failed for ", offer)

		c.conn.Close()
		return err
	}

	m = <-c.sync
	if !m.IsSetup() || !bytes.Equal(m.Payload(), []byte("ok")) {
		log.Println("the version was not ok with the remote for ", offer)
		c.conn.Close()

		return errors.New("Setup failed - wrong version of protocol")
	}

	return nil
}

func (c *connection) sendOffer(offer *msg.Msg) error {
	_, err := c.conn.Write(offer.Bytes())
	if err != nil {
		c.conn.Close()
		return err
	}

	m := <-c.sync
	if !m.IsSetup() || m.PayloadBuf().String() != "4" {
		log.Println("We didn't get versionCheck")
		log.Println("We got", m)
		c.conn.Close()
		return errors.New("VersionCheck failed")
	}

	m = msg.NewMsg(bytes.NewBufferString("ok"), msg.SETUP)

	_, err = c.conn.Write(m.Bytes())
	if err != nil {
		return err
	}

	return nil
}

func (c *connection) Close() error {
	return c.conn.Close()
}

func (c *connection) sendversionCheck() error {
	m := msg.NewMsg(
		bytes.NewBuffer([]byte{'4'}),
		msg.SETUP)
	_, err := c.conn.Write(m.Bytes())
	if err != nil {
		return err
	}
	return nil
}

type secondaryConnection struct {
	*connection
	parent *connection
}

func newSecondaryConnection(connection *connection, parent *connection) (*secondaryConnection, error) {
	c := new(secondaryConnection)
	c.connection = connection
	c.parent = parent
	return c, nil
}
