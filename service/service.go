package service

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/adrianosela/rdtp"
	"github.com/adrianosela/rdtp/network"
	"github.com/adrianosela/rdtp/packet"
	"github.com/adrianosela/rdtp/service/ports/controller"
	"github.com/pkg/errors"
)

// Service is an abstraction of the rdtp service
type Service struct {
	ports   controller.Controller
	network network.Network
}

// NewService returns an rdtp service instance
func NewService() (*Service, error) {
	ipv4Network, err := network.NewIPv4()
	if err != nil {
		return nil, errors.Wrap(err, "could not acquire network")
	}
	return &Service{
		ports:   controller.NewMemoryController(),
		network: ipv4Network,
	}, nil
}

// Run runs the rdtp service
func (s *Service) Run() error {
	// receive all rdtp packets passed on by the network
	// and forward them to the corresponding socket
	s.network.StartReceiver(func(p *packet.Packet) error {
		if err := s.ports.Deliver(p); err != nil {
			// TODO: send error message outbound
			return errors.Wrap(err, "could not deliver packet to rdtp socket")
		}
		return nil
	})

	clients, err := safeUnixListener(rdtp.DefaultRDTPServiceAddr)
	if err != nil {
		return errors.Wrap(err, "could not start system's rdtp client listener")
	}
	log.Println("[rdtp] service running")

	for {
		conn, err := clients.Accept()
		if err != nil {
			// connection closed by a signal
			if err.Error() == fmt.Sprintf("accept unix %s: use of closed network connection", rdtp.DefaultRDTPServiceAddr) {
				log.Printf("[rdtp] service stopped\n")
				return nil
			}
			// propagate any other error upstream
			log.Println(errors.Wrap(err, "could not accept rdtp client connection"))
			return err
		}
		go s.handleClientMessage(conn)
	}
}

func safeUnixListener(unixAddr string) (net.Listener, error) {
	l, err := net.Listen("unix", unixAddr)
	if err != nil {
		return nil, errors.Wrap(err, "could not listen on default RTDP service address")
	}

	// Unix sockets must be unlink()ed before being reused again.
	// Handle common process-killing signals so we can gracefully shut down:
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, os.Kill, syscall.SIGTERM)
	go func(c chan os.Signal) {
		sig := <-c
		log.Printf("[rdtp] received signal <%s>, shutting down...\n", sig)
		l.Close()
	}(sigs)

	return l, nil
}
