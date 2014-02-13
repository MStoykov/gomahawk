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
	id        string
	lastPing  time.Time
	pingTimer *time.Timer
	idbConn   *dBConn
	odbConn   *dBConn
	dbsyncKey string
	cm        *connectionManager
}

func (c *controlConnection) Name() string {
	return c.id
}

func (c *controlConnection) UUID() string {
	return c.id
}

func (t *controlConnection) RequestStreamConnection(id int64) (StreamConnection, error) {
	offerMsg := msg.NewFileRequestOffer(id, t.id)
	conn, err := t.cm.copyConnection(t.connection)
	if err != nil {
		return nil, err
	}

	sc, err := openNewStreamConnection(id, conn, t, offerMsg)

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
	_, err := msg.MakePingMsg().WriteTo(c.conn)
	return err
}

func newControlConnection(g Gomahawk, cm *connectionManager, conn *connection, id string) (*controlConnection, error) {
	c := new(controlConnection)
	c.cm = cm

	c.g = g
	c.connection = conn

	c.id = id

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
	sc, err := newSecondaryConnection(conn, c)
	if err != nil {
		return err
	}

	dbConn, err := newDBConn(sc)

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
	_, err := offer.WriteTo(c.conn)
	return err
}

func (c *controlConnection) handleMsg(m *msg.Msg) error {
	if m.IsPing() {
		c.lastPing = time.Now()
		return nil
	}

	dbsyncOffer, err := msg.ParseDBSyncOffer(m)
	if err == nil {
		log.Println("we got offer for DBSYNC ", dbsyncOffer)
		conn, err := c.cm.copyConnection(c.connection)
		if err != nil {
			return err
		}

		c.odbConn, err = openNewDBConn(dbsyncOffer, conn, c.id)
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
