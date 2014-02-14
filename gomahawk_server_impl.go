package gomahawk

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"sync"

	"github.com/MStoykov/gomahawk/udp"
	gouuid "github.com/nu7hatch/gouuid"
)

type gomahawkServerImpl struct {
	gomahawk          Gomahawk
	uuid              *gouuid.UUID
	advertisers       []udp.Advertiser
	listeners         []*tcpListener
	tomahawks         []*tomahawkImpl
	mutex             sync.Mutex
	connectionManager *connectionManager
}

// returns new instance of Gomahawk
func NewGomahawkServer(gomahawk Gomahawk) (result GomahawkServer, err error) {
	g := new(gomahawkServerImpl)
	g.gomahawk = gomahawk
	g.connectionManager = newConnectionManager()
	err = g.InitUUID()
	if err == nil {
		result = g
	}

	return
}

func (g *gomahawkServerImpl) Name() string {
	return g.gomahawk.Name()
}

func (g *gomahawkServerImpl) String() string {
	return g.Name()
}

func (g *gomahawkServerImpl) InitUUID() error {
	uuid, err := gouuid.NewV5(gouuid.NamespaceURL, []byte(g.Name()))
	if err != nil {
		return err
	}
	log.Printf("%s.uuid = %s", g, uuid)
	g.uuid = uuid

	return nil
}

func (g *gomahawkServerImpl) newConnectionCallback(tcpConn *net.TCPConn) (err error) {
	if !g.gomahawk.ConnectionIsRequested(tcpConn.RemoteAddr()) {
		tcpConn.Close() // We don't wanna speak with this one
		return nil
	}

	g.mutex.Lock()
	defer g.mutex.Unlock()

	conn := g.connectionManager.registerConnection(tcpConn)

	if err = conn.receiveOffer(); err != nil {
		tcpConn.Close()
		return err

	}
	if conn.NodeId == "" && conn.ControlId != "" {
		for _, tomahawk := range g.tomahawks {
			if tomahawk.NodeId == conn.ControlId {
				tomahawk.AddConnection(conn)
				return nil
			}
		}
		return fmt.Errorf("A connection that looks like secondary but we don't know about it's Control Connection %s", conn.OfferMsg)
	} else if conn.NodeId != "" && conn.ControlId == "" {
		cc, err := newControlConnection(g.gomahawk, g.connectionManager, conn)
		if err != nil {
			conn.Close()
			return err
		}

		newTomahawk, err := newTomahawk(cc)

		if err != nil {
			return err
		}

		g.tomahawks = append(g.tomahawks, newTomahawk)
	} else {
		return fmt.Errorf("Unhandled connection that is neither control nor secondary %s", conn.OfferMsg)
	}

	return nil
}

func (g *gomahawkServerImpl) ListenTo(ip net.IP, port int) (err error) {
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

func (g *gomahawkServerImpl) Start() (err error) {
	log.Printf("%s.Start()", g)
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

	for index, tomahawk := range g.tomahawks {
		result[index] = tomahawk
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
