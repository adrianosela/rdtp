package service

import (
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/adrianosela/rdtp"
	"github.com/adrianosela/rdtp/ipv4"
	"github.com/adrianosela/rdtp/packet"
	"github.com/adrianosela/rdtp/socket"
	"github.com/pkg/errors"
)

// Service is an abstraction of the rdtp service
type Service struct {
	favIP net.IP

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

	// FIXME: worry about listeners later

	rhost, rport, err := extractAddress(c)
	if err != nil {
		return errors.Wrap(err, "could not extract destination")
	}

	// FIXME: allocate output port here
	lhost, lport := getOutboundIP(), uint16(0) // rand.Intn(int(rdtp.MaxPort)-1)+1

	sck, err := socket.NewSocket(socket.Config{
		LocalAddr:          &rdtp.Addr{Host: lhost, Port: lport},
		RemoteAddr:         &rdtp.Addr{Host: rhost, Port: rport},
		ToApplicationLayer: c,
		ToController:       s.netLayer.Send, /* FIXME */
	})
	if err != nil {
		return errors.Wrap(err, "could not get socket for user")
	}

	if err = s.sckmgr.Put(sck); err != nil {
		return errors.Wrap(err, "could not attach socket to socket manager")
	}
	log.Printf("%s [attached]", sck.ID())

	defer func() {
		s.sckmgr.Evict(sck.ID())
		log.Printf("%s [evicted]", sck.ID())
	}()

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

// extractAddress extracts the IPv4 address and rdtp
// port of the destination address
func extractAddress(c net.Conn) (string, uint16, error) {
	buf := make([]byte, 15) // maxIP = 255.255.255.255 (15 chars)

	n, err := c.Read(buf)
	if err != nil {
		return "", uint16(0), errors.Wrap(err, "could not read from conn")
	}

	address := string(buf[:n])

	var host string
	var port uint16

	if !strings.Contains(address, ":") {
		host = address
		port = rdtp.DiscoveryPort
	} else {
		hostStr, portStr, err := net.SplitHostPort(address)
		if err != nil {
			return "", uint16(0), errors.Wrap(err, "could not split host from port")
		}
		host = hostStr

		if portStr == "" {
			port = 0
		} else {
			p64, err := strconv.ParseUint(portStr, 10, 16)
			if err != nil {
				return "", uint16(0), errors.Wrap(err, "could parse port number")
			}
			port = uint16(p64)
		}
	}

	// FIXME: DNS lookup if not IP
	return host, port, nil
}

// get preferred outbound ip of this machine
func getOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.String()
}
