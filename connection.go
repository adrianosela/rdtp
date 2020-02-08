package rdtp

import (
	"syscall"

	"github.com/adrianosela/rdtp/packet"
	"github.com/pkg/errors"
)

// Conn is an RDTP connection
type Conn struct {
	lPort uint16 // local port
	rPort uint16 // remote port

	rAddr *syscall.SockaddrInet4 // remote address

	sockFD int // socket file descriptor
}

// Dial establishes an RDTP connection with a remote IP host
func Dial(ip [4]byte) (*Conn, error) {
	// get raw network socket (AF_INET = IPv4) to send messages on
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, IPProtoRDTP)
	if err != nil {
		return nil, errors.Wrap(err, "could not get raw network socket")
	}

	// TODO: get port from RDTP controller
	// if no rdtp controller, return error
	localPort := uint16(10)

	c := &Conn{
		lPort:  localPort,
		rPort:  DiscoveryPort, // until syn ack
		rAddr:  &syscall.SockaddrInet4{Addr: ip},
		sockFD: fd,
	}

	c.syn()
	// block and wait for SYN ACK
	// repond with ACK

	return c, nil
}

func (c *Conn) syn() error {
	p, err := packet.NewPacket(c.lPort, DiscoveryPort, nil)
	if err != nil {
		return errors.Wrap(err, "could not build rdtp SYN packet for sending")
	}
	if err = syscall.Sendto(c.sockFD, p.Serialize(), 0, c.rAddr); err != nil {
		errors.Wrap(err, "could not send syn to network socket")
	}
	return nil
}
