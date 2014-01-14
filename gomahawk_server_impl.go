package gomahawk

import (
	"log"
	"net"
	"strconv"
	"sync"
	"time"

	gouuid "github.com/nu7hatch/gouuid"
)

type gomahawkServerImpl struct {
	g            Gomahawk
	uuid         *gouuid.UUID
	advertPeriod time.Duration
	advertTimer  *time.Timer
	connection   net.Conn
	addresses    []addr
	tomahawks    []*tomahawkImpl
	mutex        sync.Mutex
}

type addr struct {
	ip   net.IP
	port string
}

func (a addr) TCPAddr() (*net.TCPAddr, error) {
	portI, err := strconv.Atoi(a.port)
	if err != nil {
		return nil, err
	}

	return &net.TCPAddr{
		IP:   a.ip,
		Port: portI,
	}, nil
}

func (a addr) UDPAddr() (*net.UDPAddr, error) {
	portI, err := strconv.Atoi(a.port)
	if err != nil {
		return nil, err
	}

	return &net.UDPAddr{
		IP:   a.ip,
		Port: portI,
	}, nil
}

// returns new instance of Gomahawk
func NewGomahawkServer(g Gomahawk) (result GomahawkServer, err error) {
	gs := new(gomahawkServerImpl)
	gs.g = g
	err = gs.InitUUID()
	if err == nil {
		result = gs
	}

	//gs.connManager, err = gomahawk_net.NewConnectionManager()

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
	c, err := newConnection(conn)
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
		cc, err := newControlConnection(g.g, c, c.NodeId)
		if err != nil {
			c.Close()
			return err
		}
		newTomahawk, err := newTomahawk(cc)

		if err != nil {
			return err
		}

		g.tomahawks = append(g.tomahawks, newTomahawk) // Race Condition

	}

	return nil
}

func (g *gomahawkServerImpl) ListenTo(ip net.IP, port string) error {
	log.Printf("%s.ListenTo(%s, %s)", g, ip, port)
	newAddr := addr{ip, port}
	g.addresses = append(g.addresses, newAddr) // RACE CONDITION
	tcpAddr, err := newAddr.TCPAddr()
	if err != nil {
		return err
	}
	listener, err := newTCPListener(tcpAddr, g.newConnectionCallback)
	if err != nil {
		return err
	}
	log.Println("new listener :", listener)

	return nil
}

func (g *gomahawkServerImpl) AdvertEvery(period time.Duration) {
	g.advertPeriod = period
	if g.advertTimer != nil {
		g.advertTimer.Reset(g.advertPeriod)
	}
}

func (g *gomahawkServerImpl) Start() error {
	log.Printf("%s.Start()", g)

	g.advertTimer = time.AfterFunc(g.advertPeriod, func() {
		defer g.advertTimer.Reset(g.advertPeriod)
		err := g.AdvertNow()
		if err != nil {
			log.Println(err)
		}
	})
	if g.advertPeriod > 0 {
		g.AdvertNow()
	}

	return nil
}

func (g *gomahawkServerImpl) GetTomahawks() []Tomahawk {
	var result []Tomahawk = make([]Tomahawk, len(g.tomahawks))

	return result
}

func (g *gomahawkServerImpl) AdvertNow() error {
	//log.Printf("%s.AdvertNow()", g)
	for _, addr := range g.addresses {
		udpAddr, err := addr.UDPAddr()
		if err != nil {
			log.Println(err)
		} else {
			advert(udpAddr, g.uuid)
		}
	}

	return nil
}
