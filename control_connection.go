package gomahawk

import (
	"log"
	"time"

	msg "github.com/MStoykov/gomahawk/msg"
	gouuid "github.com/nu7hatch/gouuid"
)

type controlConnection struct {
	g Gomahawk
	*connection
	lastPing  time.Time
	pingTimer *time.Timer
	idbConn   *dBConn
	odbConn   *dBConn
	dbsyncKey string
	cm        *connectionManager
}

func (c *controlConnection) Name() string {
	return c.NodeId
}

func (c *controlConnection) UUID() string {
	return c.NodeId
}

func (c *controlConnection) RequestStreamConnection(id int64) (StreamConnection, error) {
	offerMsg := msg.NewFileRequestOffer(id, c.NodeId)
	conn, err := c.cm.copyConnection(c.connection)
	if err != nil {
		return nil, err
	}

	sc, err := openNewStreamConnection(id, conn, c, offerMsg)

	if err != nil {
		return nil, err
	}
	return sc, nil
}

func (c *controlConnection) RequestDBConnection(DBConnection) error {
	return c.sendDBSyncOffer()
}

func (c *controlConnection) LastPing() time.Time {
	return c.lastPing

}
func (c *controlConnection) sendPing() error {
	return c.WriteMsg(msg.MakePingMsg())
}

func newControlConnection(g Gomahawk, cm *connectionManager, conn *connection) (*controlConnection, error) {
	c := new(controlConnection)
	c.cm = cm

	c.g = g
	c.connection = conn

	c.msgHandler = c.handleMsg

	c.StartHandelingMessages()
	c.setupPingTimer()

	go c.sendDBSyncOffer()

	return c, nil
}

func (c *controlConnection) setupPingTimer() {
	pingTime := 5 * time.Second
	c.pingTimer = time.AfterFunc(pingTime, func() {
		defer c.pingTimer.Reset(pingTime)
		err := c.sendPing()
		if err != nil {
			log.Printf("stopping because of error : %s", err)
			return
		}
	})
}

func (c *controlConnection) AddConnection(conn *connection) error {
	dbConn, err := newDBConn(newSecondaryConnection(conn, c))

	if err != nil {
		return err
	}

	c.idbConn = dbConn

	twin, err := c.g.NewDBConnection(c, c.idbConn)
	if err != nil {
		c.idbConn.Close()
		c.idbConn = nil
		return err
	}
	c.idbConn.dbconn = twin
	c.idbConn.StartHandelingMessages()

	return nil
}

func (c *controlConnection) sendDBSyncOffer() error {
	uuid, _ := gouuid.NewV4()
	c.dbsyncKey = uuid.String()

	offer := msg.NewDBSyncOfferMsg(c.dbsyncKey)
	return c.WriteMsg(offer)
}

func (c *controlConnection) handleMsg(m *msg.Msg) error {
	if m.IsPing() {
		c.lastPing = time.Now()
		return nil
	}

	dbsyncOffer, err := msg.ParseDBSyncOffer(m)
	if err == nil {
		conn, err := c.cm.copyConnection(c.connection)
		if err != nil {
			return err
		}

		c.odbConn, err = openNewDBConn(dbsyncOffer, conn, c.NodeId)
		if err != nil {
			return err
		}

		dbConn, err := c.g.NewDBConnectionRequested(c, c.odbConn)
		if err != nil {
			log.Println(err)
			return err
		}
		log.Println(dbConn)
		// DO something more incase of error
		return err
	}

	log.Println("unhandled m :", m)
	// parse
	return nil
}
