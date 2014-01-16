package gomahawk

import (
	"log"
	"net"
	"strconv"
	"sync"

	"github.com/MStoykov/gomahawk/udp"
	gouuid "github.com/nu7hatch/gouuid"
)

type gomahawkServerImpl struct {
	g           Gomahawk
	uuid        *gouuid.UUID
	advertisers []udp.Advertiser
	listeners   []*tcpListener
	tomahawks   []*tomahawkImpl
	mutex       sync.Mutex
	cm          *connectionManager
}

// returns new instance of Gomahawk
func NewGomahawkServer(g Gomahawk) (result GomahawkServer, err error) {
	gs := new(gomahawkServerImpl)
	gs.g = g
	err = gs.InitUUID()
	if err == nil {
		result = gs
	}

	gs.cm = newConnectionManager()

	if err != nil {
		log.Printf("error whie initializing ConnectionManager : %s", err)
	}
	return
}

func (g *gomahawkServerImpl) Name() string {
	return g.g.Name()
}

func (g *gomahawkServerImpl) String() string {
	return g.Name()
}

func (g *gomahawkServerImpl) InitUUID() error {
	log.Printf("%s.InitUUID()", g)
	uuid, err := gouuid.NewV5(gouuid.NamespaceURL, []byte(g.Name()))
	if err != nil {
		return err
	}
	log.Printf("%s.uuid = %s", g, uuid)
	g.uuid = uuid

	return nil
}

func (g *gomahawkServerImpl) newConnectionCallback(conn *net.TCPConn) error {
	log.Println("newConnectionCallback(", conn, ")")
	if !g.g.ConnectionIsRequested(conn.RemoteAddr()) {
		conn.Close() // We don't wanna speak with this one
		return nil
	}

	g.mutex.Lock()
	defer g.mutex.Unlock()

	c, err := g.cm.newConnection(conn)

	if err != nil {
		conn.Close()
		return nil

	}
	err = c.receiveOffer()
	if err != nil {
		conn.Close()
		return nil

	}
	if c.NodeId == "" && c.ControlId != "" {
		for _, tomahawk := range g.tomahawks {
			if tomahawk.NodeId == c.ControlId {
				log.Printf("adding conn [%s] to Tomahawk[%s]", c, tomahawk)
				tomahawk.AddConnection(c)
				return nil
			}
		}
	} else if c.NodeId != "" && c.ControlId == "" {
		cc, err := newControlConnection(g.g, g.cm, c, c.NodeId)
		if err != nil {
			c.Close()
			return err
		}
		newTomahawk, err := newTomahawk(cc)

		if err != nil {
			return err
		}

		g.tomahawks = append(g.tomahawks, newTomahawk)
	}

	return nil
}

func (g *gomahawkServerImpl) ListenTo(ip net.IP, port int) error {
	log.Printf("%s.ListenTo(%s, %s)", g, ip, port)
	listener, err := newTCPListener(ip, port, g.newConnectionCallback)
	if err != nil {
		return err
	}
	log.Println("new listener :", listener)

	advertiser, err := udp.NewAdvertiser(ip, g.uuid.String(), g.Name(), strconv.Itoa(port))
	if err != nil {
		listener.Close()
		return err
	}
	g.advertisers = append(g.advertisers, advertiser)
	g.listeners = append(g.listeners, listener)

	return nil
}

func (g *gomahawkServerImpl) AdvertiseEvery(seconds int) {
	for _, advertiser := range g.advertisers {
		advertiser.AdvertiseEvery(seconds)
	}
}

func (g *gomahawkServerImpl) Start() error {
	log.Printf("%s.Start()", g)
	var err error
	for _, advertiser := range g.advertisers {
		if err = advertiser.Start(); err != nil {
			return err
		}
	}

	for _, tcpListener := range g.listeners {
		if err = tcpListener.Start(); err != nil {
			return err
		}
	}
	return nil
}

func (g *gomahawkServerImpl) GetTomahawks() []Tomahawk {
	var result []Tomahawk = make([]Tomahawk, len(g.tomahawks))

	for i, tomahawk := range g.tomahawks {
		result[i] = tomahawk
	}
	return result
}

func (g *gomahawkServerImpl) Advertise() (err error) {
	for _, advertiser := range g.advertisers {
		if err = advertiser.Advertise(); err != nil {
			return err
		}
	}

	return nil
}
