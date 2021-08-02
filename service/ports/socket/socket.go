package socket

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/adrianosela/rdtp"
	"github.com/adrianosela/rdtp/network"
	"github.com/adrianosela/rdtp/packet"
	"github.com/adrianosela/rdtp/packet/factory"
	"github.com/pkg/errors"
)

const (
	inboundPacketChannelSize = 100
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

	// packetizes and forwards to network layer
	packetizer *factory.PacketFactory

	// packets received at the network
	// are ultimately delivered in this
	// channel to be read by the socket
	// and be written as messages to
	// the application layer
	inbound chan *packet.Packet

	// used to notify socket of shutdown
	shutdown chan bool
}

// Config is the necessary configuration to initialize a socket
type Config struct {
	LocalAddr  *rdtp.Addr // local rdtp address
	RemoteAddr *rdtp.Addr // remote rdtp address

	// connection to app layer
	Application net.Conn

	// connection to network layer
	Network network.Network
}

// New is the socket constructor
func New(c Config) (*Socket, error) {
	if c.LocalAddr == nil || net.ParseIP(c.LocalAddr.Host) == nil {
		return nil, errors.New("invalid local address")
	}
	if c.RemoteAddr == nil || net.ParseIP(c.LocalAddr.Host) == nil {
		return nil, errors.New("remote address cannot be nil")
	}
	if c.Application == nil {
		return nil, errors.New("connection to application layer cannot be nil")
	}
	if c.Network == nil {
		return nil, errors.New("connection to network layer cannot be nil")
	}

	toNetwork := func(p *packet.Packet) error {
		p.SetSourceIPv4(net.ParseIP(c.LocalAddr.Host))
		p.SetDestinationIPv4(net.ParseIP(c.RemoteAddr.Host))
		return c.Network.Send(p)
	}

	return &Socket{
		lAddr:       c.LocalAddr,
		rAddr:       c.RemoteAddr,
		application: c.Application,
		packetizer: factory.DefaultPacketFactory(
			net.ParseIP(c.LocalAddr.Host),
			net.ParseIP(c.RemoteAddr.Host),
			uint16(c.LocalAddr.Port),
			uint16(c.RemoteAddr.Port),
			toNetwork),
		inbound:  make(chan *packet.Packet, inboundPacketChannelSize),
		shutdown: make(chan bool, 1),
	}, nil
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
	// close conn to application layer
	s.application.Close()
}

// Deliver delivers a packet to a socket's inbound packet channel
func (s *Socket) Deliver(p *packet.Packet) {
	s.inbound <- p
}

// Run kicks-off socket processes
func (s *Socket) Run() {
	done := make(chan bool, 1)

	go s.receive(done)
	go s.transmit()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case <-sigs:
		case <-s.shutdown:
			done <- true
			close(done)
			// FIXME: proper FIN handshake -- for now just sending FIN and returning
			s.packetizer.SendControlPacket(false, false, true, false)
			return
		}
	}
}

func (s *Socket) receive(done chan bool) {
	for {
		select {
		case <-done:
			close(s.inbound) // TODO: move to after FIN handshake
			return
		case p := <-s.inbound:
			s.rxBytes += uint32(p.Length)  // stats
			s.application.Write(p.Payload) // pass packet to application layer
		}
	}
}

func (s *Socket) transmit() {
	buf := make([]byte, 1500)
	for {
		n, err := s.application.Read(buf)
		if err != nil {
			if err == io.EOF {
				s.shutdown <- true
				close(s.shutdown)
				return
			}
			continue
		}

		n, err = s.packetizer.PackAndForwardMessage(buf[:n])
		if err != nil {
			log.Printf("[rdtp socket %s] Error packetizing and forwarding message: %s", s.ID(), err)
			return
		}

		s.txBytes += uint32(n) // stats
	}
}
