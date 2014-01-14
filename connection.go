package gomahawk

import (
	"bytes"
	"errors"
	msg "github.com/MStoykov/gomahawk/msg"
	"log"
	"net"
)

type connection struct {
	conn      *net.TCPConn
	processor *msg.Processor
	sync      <-chan *msg.Msg
	*msg.OfferMsg
}

func (c *connection) setupProcessor() {
	sync := make(chan *msg.Msg)
	c.sync = sync
	c.processor = msg.NewProcessor(c.conn, sync)
}

func newConnection(conn *net.TCPConn) (*connection, error) {
	c := new(connection)

	c.conn = conn

	c.setupProcessor()

	m := <-c.sync
	offer, err := msg.ParseOffer(m)
	if err != nil {
		return nil, err
	}
	c.OfferMsg = offer
	if err := c.sendversionCheck(); err != nil {
		log.Println("versionCheck sending failed for ", offer)

		c.conn.Close()
		return nil, err
	}

	m = <-c.sync
	if !m.IsSetup() || !bytes.Equal(m.Payload(), []byte("ok")) {
		log.Println ("the version was not ok with the remote for ", offer)
		c.conn.Close()

		return nil, errors.New("Setup failed - wrong version of protocol")
	}


	return c, nil
}

func openNewConnection(local net.Addr, remote net.Addr, offer *msg.Msg) (*connection, error) {
	tcpLAddr, err := net.ResolveTCPAddr("tcp", local.String())
	if err != nil {
		log.Println("error while resolving local tcp addr", err)
		return nil, err
	}
	tcpLAddr.Port = 0 // we need to get new port and we don't care which one it is :)
	tcpRAddr, err := net.ResolveTCPAddr("tcp", remote.String())
	if err != nil {
		log.Println("error while resolving remote tcp addr", err)
		return nil, err
	}
	tcpRAddr.Port = 50210 // hardcoded

	conn, err := net.DialTCP("tcp", tcpLAddr, tcpRAddr)
	if err != nil {
		log.Println("couldn't dial remote ", remote, "from local", local)
		return nil, err
	}

	c := new(connection)
	c.conn = conn
	c.setupProcessor()
	
	_, err = c.conn.Write(offer.Bytes())
	if err != nil {
		c.conn.Close()
		return nil, err
	}

	m := <-c.sync
	if !m.IsSetup() || m.PayloadBuf().String() != "4" {
		log.Println("We didn't get versionCheck")
		log.Println("We got", m)
		c.conn.Close()
		return nil, errors.New("VersionCheck failed")
	}

	m = msg.NewMsg(bytes.NewBufferString("ok"), msg.SETUP)

	_, err = c.conn.Write(m.Bytes())
	if err != nil {
		return nil, err
	}

	return c, nil
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
