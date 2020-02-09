package service

import (
	"fmt"
	"log"
	"net"
	"os"
	"syscall"

	"github.com/adrianosela/rdtp"
	"github.com/adrianosela/rdtp/packet"
	"github.com/adrianosela/rdtp/ports"
	"github.com/adrianosela/rdtp/ports/filesystem"
	"github.com/pkg/errors"
)

// Service represents the RDTP service which is in charge of two tasks:
// - handing out RDTP connections to service clients on "this" host
// - listening for RDTP packets over IP and multiplexing to their rdtp client
type Service struct {
	unixSock string
	ports    ports.Manager
	mux      map[uint16]net.Conn // port number to client conn
}

// NewService returns the default RDTP service
func NewService() (*Service, error) {
	mgr, err := filesystem.NewFSManager("")
	if err != nil {
		return nil, errors.Wrap(err, "could not init file system ports manager")
	}
	return &Service{
		unixSock: "/tmp/rdtp2.sock",
		ports:    mgr,
	}, nil
}

// Start starts the RDTP service
func (s *Service) Start() error {
	// l, err := net.Listen("unix", s.unixSock)
	// if err != nil {
	// 	return errors.Wrap(err, "could not listen on unix socket")
	// }
	// defer l.Close()

	go s.listenRDTP()
	for {
		// conn, err := l.Accept()
		// if err != nil {
		// 	log.Println(errors.Wrap(err, "could not accept rdtp client connection"))
		// 	continue
		// }
		// go s.serveClient(conn)
	}
}

func (s *Service) listenRDTP() {
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, rdtp.IPProtoRDTP)
	if err != nil {
		log.Fatal(errors.Wrap(err, "could not get raw network socket"))
	}
	f := os.NewFile(uintptr(fd), fmt.Sprintf("fd %d", fd))

	fmt.Println("listening for RDTP on all local network interfaces")
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

		if !rdtpPacket.Check() {
			log.Println("failed checksum, packet dropped")
			continue
		}

		if err = s.MultiplexPacket(rdtpPacket); err != nil {
			log.Println(errors.Wrap(err, "could not multiplex rdtp packet"))
			continue
		}
	}
}

func (s *Service) serveClient(c net.Conn) {
	p, err := s.ports.AllocateAny()
	if err != nil {
		log.Println(errors.Wrap(err, "[RDTP] could not allocate port for client"))
		return
	}

	s.mux[p] = c // attach client connection to MUX

	// TODO
}
