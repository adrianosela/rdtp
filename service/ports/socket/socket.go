package socket

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/adrianosela/rdtp"
	"github.com/adrianosela/rdtp/network"
	"github.com/adrianosela/rdtp/packet"
	"github.com/adrianosela/rdtp/packet/factory"
	"github.com/pkg/errors"
)

const (
	flagFmt = "{SYN[%t] ACK[%t] FIN[%t] ERR[%t]}"

	packetChannelSize = 100
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
		inbound:  make(chan *packet.Packet, packetChannelSize),
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
	// shutdown reader/writer threads
	s.shutdown <- true
	close(s.shutdown)
	// close conn to application layer
	s.application.Close()
}

// receiveControlPacket blocks until the a packet is received (or timeout)
func receiveControlPacket(in chan *packet.Packet, syn, ack, fin, err bool, timeout time.Duration) error {
	for {
		select {
		case p := <-in:
			if syn != p.IsSYN() || ack != p.IsACK() || fin != p.IsFIN() || err != p.IsERR() {
				return fmt.Errorf(
					"expected packet with flags %s but got %s",
					fmt.Sprintf(flagFmt, syn, ack, fin, err),
					fmt.Sprintf(flagFmt, p.IsSYN(), p.IsACK(), p.IsFIN(), p.IsERR()))
			}
			return nil
		case <-time.After(timeout):
			return errors.New("operation timed out")
		}
	}
}

// Dial sends a SYN, waits for a SYN ACK, and sends an ACK
func (s *Socket) Dial() error {
	// send SYN
	if err := s.packetizer.SendControlPacket(true, false, false, false); err != nil {
		return errors.Wrap(err, "handshake failed when sending SYN")
	}
	// wait for SYN ACK
	if err := receiveControlPacket(s.inbound, true, true, false, false, time.Second*1); err != nil {
		return errors.Wrap(err, "handshake failed when waiting for SYN ACK")
	}
	// send ACK
	if err := s.packetizer.SendControlPacket(false, true, false, false); err != nil {
		return errors.Wrap(err, "handshake failed when sending ACK")
	}
	return nil
}

// Accept sends a SYN ACK and waits for an ACK
func (s *Socket) Accept() error {
	// send SYN ACK
	if err := s.packetizer.SendControlPacket(true, true, false, false); err != nil {
		return errors.Wrap(err, "handshake failed when sending SYN ACK")
	}
	// wait for ACK
	if err := receiveControlPacket(s.inbound, false, true, false, false, time.Second*1); err != nil {
		return errors.Wrap(err, "handshake failed when waiting for ACK")
	}
	return nil
}

// Deliver delivers a packet to a socket's inbound packet channel
func (s *Socket) Deliver(p *packet.Packet) {
	s.inbound <- p
}

// Run kicks-off socket processes
func (s *Socket) Run() error {
	rxdone := make(chan bool, 1)
	txdone := make(chan bool, 1)

	go s.receive(rxdone)
	go s.transmit(txdone)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case <-sigs:
		case <-s.shutdown:
			txdone <- true
			rxdone <- true
			close(txdone)
			close(rxdone)
			// FIXME: proper FIN handshake -- for now just sending FIN and returning
			s.packetizer.SendControlPacket(false, false, true, false)
			return nil
		}
	}
}

func (s *Socket) receive(done chan bool) {
	for {
		select {
		case <-done:
			close(s.inbound)
			return
		case p := <-s.inbound:
			s.rxBytes += uint32(p.Length)  // stats
			s.application.Write(p.Payload) // pass packet to application layer
		}
	}
}

func (s *Socket) transmit(done chan bool) {
	buf := make([]byte, 1500)
	for {
		select {
		case <-done:
			return
		default:
			n, err := s.application.Read(buf)
			if err != nil {
				if err == io.EOF {
					s.shutdown <- true // client closed conn, shutdown socket
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
}
