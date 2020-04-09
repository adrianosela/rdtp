package socket

import (
	"fmt"
	"log"
	"net"

	"github.com/adrianosela/rdtp"
	"github.com/adrianosela/rdtp/factory"
	"github.com/adrianosela/rdtp/packet"
	"github.com/pkg/errors"
)

// Socket represents a socket abstraction and carries all
// necessary info and statistics about the socket
type Socket struct {
	lAddr *rdtp.Addr // local rdtp address
	rAddr *rdtp.Addr // remote rdtp address

	txBytes uint32 // current sequence number
	rxBytes uint32 // current ack number

	// connection to app layer
	application net.Conn

	// messages are read from application
	// and sent to the packet factory for
	// formatting and to then be sent out
	outbound *factory.PacketFactory

	// packets received at the network
	// are ultimately delivered in this
	// channel to be read by the socket
	// and be written as messages to
	// the application layer
	inbound chan *packet.Packet
}

// Config is the necessary configuration to initialize a socket
type Config struct {
	LocalAddr  *rdtp.Addr // local rdtp address
	RemoteAddr *rdtp.Addr // remote rdtp address

	ToApplicationLayer net.Conn
	ToController       func(p *packet.Packet) error
}

// NewSocket returns a newly allocated socket
func NewSocket(c Config) (*Socket, error) {
	if c.LocalAddr == nil {
		return nil, errors.New("local address cannot be nil")
	}
	if c.RemoteAddr == nil {
		return nil, errors.New("remote address cannot be nil")
	}
	if c.ToApplicationLayer == nil {
		return nil, errors.New("connection to application layer cannot be nil")
	}
	if c.ToController == nil {
		return nil, errors.New("connection to controller cannot be nil")
	}

	lp, rp := uint16(c.LocalAddr.Port), uint16(c.RemoteAddr.Port)

	outbound := func(p *packet.Packet) error {
		local := net.ParseIP(c.LocalAddr.Host)
		rmte := net.ParseIP(c.RemoteAddr.Host)
		// TODO: error check the IPs somewhere...
		p.SetSourceIPv4(local)
		p.SetDestinationIPv4(rmte)
		return c.ToController(p)
	}

	pf, err := factory.New(lp, rp, outbound, packet.MaxPayloadBytes)
	if err != nil {
		return nil, errors.Wrap(err, "could not initialize new packetfactory")
	}

	return &Socket{
		lAddr:       c.LocalAddr,
		rAddr:       c.RemoteAddr,
		application: c.ToApplicationLayer,
		outbound:    pf,
		inbound:     make(chan *packet.Packet),
	}, nil
}

// Start kicks-off socket processes
func (s *Socket) Start() error {
	go s.receive()
	go s.transmit()

	for {
		/* TODO:
		handle all disconnections here
		*/
	}
}

// ID returns the of unique identifier of the socket
func (s *Socket) ID() string {
	return fmt.Sprintf("%s %s", s.lAddr.String(), s.rAddr.String())
}

// LocalAddr returns the local network address.
func (s *Socket) LocalAddr() net.Addr {
	return s.lAddr
}

// RemoteAddr returns the remote network address.
func (s *Socket) RemoteAddr() net.Addr {
	return s.rAddr
}

// Close closes a socket
func (s *Socket) Close() {
	close(s.inbound)
}

func (s *Socket) receive() {
	for {
		p := <-s.inbound

		s.rxBytes += uint32(p.Length)  // keep track of stats
		s.application.Write(p.Payload) // pass packet to application layer

		// FIXME - no end condition
	}
}

func (s *Socket) transmit() {
	buf := make([]byte, 1500)
	for {

		n, err := s.application.Read(buf)
		if err != nil {
			return // FIXME
		}

		n, err = s.outbound.Send(buf[:n])
		if err != nil {
			return // FIXME
		}
		s.txBytes += uint32(n)
		// FIXME - no end condition
	}
}
