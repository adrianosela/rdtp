package socket

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

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
	if c.LocalAddr == nil || net.ParseIP(c.LocalAddr.Host) == nil {
		return nil, errors.New("invalid local address")
	}
	if c.RemoteAddr == nil || net.ParseIP(c.LocalAddr.Host) == nil {
		return nil, errors.New("remote address cannot be nil")
	}
	if c.ToApplicationLayer == nil {
		return nil, errors.New("connection to application layer cannot be nil")
	}
	if c.ToController == nil {
		return nil, errors.New("connection to controller cannot be nil")
	}

	la, ra := net.ParseIP(c.LocalAddr.Host), net.ParseIP(c.RemoteAddr.Host)
	lp, rp := uint16(c.LocalAddr.Port), uint16(c.RemoteAddr.Port)

	outFunc := func(p *packet.Packet) error {
		p.SetSourceIPv4(la)
		p.SetDestinationIPv4(ra)
		return c.ToController(p)
	}

	outbound, err := factory.New(lp, rp, outFunc, packet.MaxPayloadBytes)
	if err != nil {
		return nil, errors.Wrap(err, "could not initialize new packetfactory")
	}

	return &Socket{
		lAddr:       c.LocalAddr,
		rAddr:       c.RemoteAddr,
		application: c.ToApplicationLayer,
		outbound:    outbound,
		inbound:     make(chan *packet.Packet, 100),
	}, nil
}

// Start kicks-off socket processes
func (s *Socket) Start() error {
	rxdone := make(chan bool, 1)
	txdone := make(chan bool, 1)

	go s.receive(rxdone)
	go s.transmit(txdone)

	// TODO: figure out how to handle disconnections
	// e.g. ensure if any conn along the way breaks,
	// close all goroutines

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	for {
		<-sigs // blocks on signal
		rxdone <- true
		txdone <- true
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

func (s *Socket) receive(done chan bool) {
	for {
		select {
		case p := <-s.inbound:
			s.rxBytes += uint32(p.Length)  // stats
			s.application.Write(p.Payload) // pass packet to application layer
		case <-done:
			break
		}
	}
}

func (s *Socket) transmit(done chan bool) {
	buf := make([]byte, 1500)
	for {
		select {
		case <-done:
			break
		default:
			n, err := s.application.Read(buf)
			if err != nil {
				return // FIXME
			}

			n, err = s.outbound.Send(buf[:n])
			if err != nil {
				return // FIXME
			}

			s.txBytes += uint32(n) // stats
		}
	}
}
