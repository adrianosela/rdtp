package service

import (
	"log"
	"net"

	"github.com/adrianosela/rdtp"
	"github.com/adrianosela/rdtp/ipv4"
	"github.com/adrianosela/rdtp/packet"
	"github.com/adrianosela/rdtp/socket"
	"github.com/pkg/errors"
)

// Service is an abstraction of the rdtp service
type Service struct {
	appLayer net.Listener
	sckmgr   *socket.Manager
	netLayer *ipv4.IPv4
}

// NewService returns an rdtp service instance
func NewService() (*Service, error) {
	network, err := ipv4.NewIPv4()
	if err != nil {
		return nil, errors.Wrap(err, "could not acquire network")
	}
	socketManager, err := socket.NewManager()
	if err != nil {
		return nil, errors.Wrap(err, "could not initialize socket manager")
	}

	applications, err := net.Listen("unix", rdtp.DefaultRDTPServiceAddr)
	if err != nil {
		return nil, errors.Wrap(err, "could not listen on default RTDP service address")
	}

	return &Service{
		appLayer: applications,
		sckmgr:   socketManager,
		netLayer: network,
	}, nil
}

// Start starts the rdtp service
func (s *Service) Start() error {
	// receive all rdtp packets passed on by the network
	// and forward them to the corresponding socket
	go s.netLayer.Receive(func(p *packet.Packet) error {
		if err := s.sckmgr.Deliver(p); err != nil {
			return errors.Wrap(err, "could not deliver packet to rdtp socket")
		}
		return nil
	})

	for {
		conn, err := s.appLayer.Accept()
		if err != nil {
			log.Println(errors.Wrap(err, "could not accept rdtp client connection"))
			continue
		}
		go s.handleUser(conn)
	}
}

func (s *Service) handleUser(c net.Conn) error {
	defer c.Close()

	// FIXME: read config here
	// FIXME: allocate port here

	sck, err := socket.NewSocket(socket.Config{
		LocalAddr:          &rdtp.Addr{Host: "127.0.0.1", Port: rdtp.Port(10)}, // FIXME
		RemoteAddr:         &rdtp.Addr{Host: "8.8.8.8", Port: rdtp.Port(10)},   // FIXME
		ToApplicationLayer: c,
		ToController:       s.netLayer.Send, /* FIXME */
	})
	if err != nil {
		return errors.Wrap(err, "could not get socket for user")
	}

	if err = s.sckmgr.Put(sck); err != nil {
		return errors.Wrap(err, "could not attach socket to socket manager")
	}
	defer s.sckmgr.Evict(sck.ID())

	if err = sck.Start(); err != nil {
		return errors.Wrap(err, "socket failure")
	}
	return nil
}
