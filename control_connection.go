package gomahawk

import (
	"bytes"
	"encoding/json"
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
	offerMsg, err := msg.NewFileRequestOffer(id, t.id)
	if err != nil {
		return nil, err
	}
	conn, err := t.cm.copyConnection(t.connection)
	if err != nil {
		return nil, err
	}

	sc, err := openNewStreamConnection(id, conn, t.connection, offerMsg)

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
	_, err := c.conn.Write(msg.MakePingMsg().Bytes())
	return err
}

func newControlConnection(g Gomahawk, cm *connectionManager, conn *connection, id string) (*controlConnection, error) {
	c := new(controlConnection)
	c.cm = cm

	c.g = g
	c.connection = conn

	c.id = id

	go func() {
		for  {
			m, err := c.processor.ReadMSG() // handle err
			err = c.handleMsg(m)
			if err != nil {
				log.Println(err)
				return
			}
		}
	}()
	defer c.sendDBSyncOffer()

	c.setupPingTimer()

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

	twin, err := c.g.NewDBConnection(c, c.idbConn)
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
