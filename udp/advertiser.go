package udp

import (
	"fmt"
	"net"
	"time"

	gouuid "github.com/nu7hatch/gouuid"
)

type advertPacket struct {
	advert string
	port   string
	uuid   string
	name   string
}

const (
	DefaultUDPPort         int = 50210
	DefaultAdvertisePeriod     = 5
)

func uuidFromHostname(hostname string) (string, error) {
	uuid, err := gouuid.NewV5(gouuid.NamespaceURL, []byte(hostname))
	if err != nil {
		return "", err
	}
	return uuid.String(), nil
}

// make new advertPacket
func newAdvertPacket(advert string, port string, uuid string, name string) *advertPacket {
	packet := advertPacket{
		advert,
		port,
		uuid,
		name,
	}
	return &packet
}

func (a *advertPacket) Bytes() []byte {
	str := fmt.Sprintf("%s:%s:%s", a.advert, a.port, a.uuid)

	if a.name != "" {
		str += ":" + a.name
	}
	return []byte(str)
}

type Advertiser interface {
	Start() error
	Stop() error
	Advertise() error
	AdvertiseEvery(seconds int)
}

type advertiser struct {
	*Broadcaster
	timer     *time.Timer
	seconds   time.Duration
	lastError error
}

func (a *advertiser) Start() error {
	a.timer = time.AfterFunc(a.seconds, func() {
		a.Advertise()
		a.timer.Reset(a.seconds)
	})

	if a.seconds > 0 {
		a.Advertise()
	}

	return nil
}

func (a *advertiser) Stop() error {
	a.timer.Stop()
	return nil
}
func (a *advertiser) Advertise() error {
	return a.Broadcast()
}
func (a *advertiser) AdvertiseEvery(seconds int) {
	a.seconds = time.Duration(seconds) * time.Second
	if a.timer != nil {
		if a.seconds > 0 {
			a.timer.Reset(a.seconds)
		} else {
			a.timer.Stop()
		}
	}
}

func NewAdvertiser(ip net.IP, uuid, name, port string) (Advertiser, error) {
	var err error

	a := new(advertiser)
	a.Broadcaster, err = NewBroadCaster(ip, DefaultUDPPort, newAdvertPacket("TOMAHAWKADVERT", port, uuid, name).Bytes())
	if err != nil {
		return nil, err
	}
	a.seconds = DefaultAdvertisePeriod

	return a, nil
}
