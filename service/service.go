package service

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/adrianosela/rdtp"
	"github.com/adrianosela/rdtp/ipv4"
	"github.com/adrianosela/rdtp/packet"
	"github.com/adrianosela/rdtp/socket"
	"github.com/pkg/errors"
)

// Service is an abstraction of the rdtp service
type Service struct {
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
	return &Service{
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

	clients, err := safeUnixListener(rdtp.DefaultRDTPServiceAddr)
	if err != nil {
		return errors.Wrap(err, "could not start system's rdtp client listener")
	}

	for {
		conn, err := clients.Accept()
		if err != nil {
			log.Println(errors.Wrap(err, "could not accept rdtp client connection"))
			continue
		}
		go s.handleClient(conn)
	}
}

func (s *Service) handleClient(c net.Conn) error {
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

func safeUnixListener(unixAddr string) (net.Listener, error) {
	l, err := net.Listen("unix", unixAddr)
	if err != nil {
		return nil, errors.Wrap(err, "could not listen on default RTDP service address")
	}

	// Unix sockets must be unlink()ed before being reused again.
	// Handle common process-killing signals so we can gracefully shut down:
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, os.Kill, syscall.SIGTERM)
	go func(c chan os.Signal) {
		sig := <-c
		log.Printf("[rdtp] signal %s: shutting down.", sig)
		l.Close()
		os.Exit(0)
	}(sigChan)

	return l, nil
}
