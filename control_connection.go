package gomahawk

import (
	"bytes"
	"encoding/json"
	"log"
	"time"

	msg "github.com/MStoykov/gomahawk/msg"
	gouuid "github.com/nu7hatch/gouuid"
)

type controlConnection struct  {
	g Gomahawk
	*connection
	id        string
	lastPing  time.Time
	pingTimer *time.Timer
	idbConn   *dBConn
	odbConn   *dBConn
	dbsyncKey string
}

func (c *controlConnection) Name() string {
	return c.id
}

func (c *controlConnection) UUID() string {
	return c.id
}

func (t *controlConnection) RequestStreamConnection(id int64)  (StreamConnection, error){
	offerMsg, err := msg.NewFileRequestOffer(id, t.id)
	if err != nil {
		return nil, err
	}
	sc, err := openNewStreamConnection(id, t.conn.LocalAddr(), t.conn.RemoteAddr(), t.connection,  offerMsg)

	log.Println("Got streamConnection")
	log.Println(sc)

	return sc, nil
}

func (c *controlConnection) RequestDBConnection(DBConnection) error {
	return c.sendDBSyncOffer()
}

func (c *controlConnection) test2132() {
}

func (c *controlConnection) LastPing()  time.Time{
	return c.lastPing

}
func (c *controlConnection) sendPing()  error {
	_, err := c.conn.Write(msg.MakePingMsg().Bytes())
	return err
}

func newControlConnection(g Gomahawk, conn *connection, id string)  (*controlConnection, error){
	c := new(controlConnection)

	c.g = g
	c.connection = conn

	c.id = id

	go func() {
		for msg := range c.sync {
			err := c.handleMsg(msg)
			if err != nil {
				log.Println(err)
				return
			}
		}
	}()
	defer c.sendDBSyncOffer()

	c.setupPingTimer()
// controlConnection

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
	sc, err := newSecondaryConnection(conn, c.connection)
	if err != nil {
		return err
	}

	dbConn, err := newDBConn(sc)

	if err != nil {
		return err
	}

	c.idbConn = dbConn

	twin , err := c.g.NewDBConnection(c, c.idbConn)
	log.Println(twin)
	if err != nil {
		dbConn.Close()
		return err 
	}



	return nil
}

func (c *controlConnection) sendDBSyncOffer() error {
	offer := new(msg.DBsyncOffer)
	offer.Method = "dbsync-offer"
	uuid, _ := gouuid.NewV4()
	offer.Key = uuid.String()
	c.dbsyncKey = offer.Key

	offerBytes, err := json.Marshal(offer)
	msg := msg.NewMsg(bytes.NewBuffer(offerBytes), msg.SETUP|msg.JSON)
	log.Println("gonna sedn offer", msg)
	_, err = c.conn.Write(msg.Bytes())
	if err != nil {
		log.Println("error while sending offer for dbconnection")
	}
	return err
}

func (c *controlConnection) handleMsg(m *msg.Msg) error {
	if m.IsPing() {
		//log.Println("pinged")
		c.lastPing = time.Now()
		return nil // MAYBE MISTAKE
	}
	dbsyncOffer, err := msg.ParseDBSyncOffer(m)
	if err == nil {
		log.Println("we got offer for DBSYNC ", dbsyncOffer)
		c.odbConn, err = openNewDBConn(dbsyncOffer, c.conn.LocalAddr(), c.conn.RemoteAddr(), c.id)
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
