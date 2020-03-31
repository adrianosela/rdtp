package service

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/adrianosela/rdtp"
	"github.com/adrianosela/rdtp/multiplexer"
	"github.com/adrianosela/rdtp/packet"
	"github.com/pkg/errors"
)

// Service represents the RDTP service which is in charge of two tasks:
// - handling RDTP connections to service clients on "this" host
// - listening for RDTP packets over IP and multiplexing to their rdtp client
type Service struct {
	unixSock string
	mux      multiplexer.Mux
}

// NewService returns the default RDTP service
func NewService() (*Service, error) {
	return &Service{
		unixSock: rdtp.DefaultRDTPServiceAddr,
		mux:      multiplexer.NewMemoryMux(),
	}, nil
}

// Start starts the RDTP service
func (s *Service) Start() error {
	l, err := net.Listen("unix", s.unixSock)
	if err != nil {
		return errors.Wrap(err, "could not listen on unix socket")
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

	go s.listenRDTP()

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println(errors.Wrap(err, "could not accept rdtp client connection"))
			continue
		}

		go s.handleUser(conn)
	}
}

func (s *Service) listenRDTP() {
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, rdtp.IPProtoRDTP)
	if err != nil {
		log.Fatal(errors.Wrap(err, "could not get raw network socket"))
	}
	f := os.NewFile(uintptr(fd), fmt.Sprintf("fd %d", fd))

	log.Printf("[rdtp] listening for RDTP packets on all local network interfaces")
	for {
		buf := make([]byte, 1500) // maximum RDTP packet

		ipDatagramSize, err := f.Read(buf)
		if err != nil {
			log.Println(errors.Wrap(err, "could not read data from network socket"))
			continue
		}

		rawIP := []byte(buf)[:ipDatagramSize]
		ihl := 4 * (rawIP[0] & byte(15))
		rawRDTP := rawIP[ihl:]

		rdtpPacket, err := packet.Deserialize(rawRDTP)
		if err != nil {
			log.Println(errors.Wrap(err, "could not deserialize rdtp packet"))
			continue
		}

		s.mux.MultiplexPacket(rdtpPacket) // note the ignored error
	}
}

func (s *Service) handleUser(c net.Conn) error {
	defer c.Close() // ensure we close conn

	// TODO: receive port number request

	p, err := s.mux.AttachAny(c)
	if err != nil {
		return errors.Wrap(err, "could not associate connection with port")
	}
	log.Printf("[rdtp] new client on port %d", p)

	// FIXME
	for {
	}
}
