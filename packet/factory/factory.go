package factory

import (
	"fmt"
	"net"

	"github.com/adrianosela/rdtp/packet"
	"github.com/pkg/errors"
)

// PacketFactory formats streams of arbitrary data onto
// chunks of a given max size, builds an rdtp packet
// with the data, and forwards it
type PacketFactory struct {
	lhost  net.IP
	rhost  net.IP
	lport  uint16
	rport  uint16
	fwFunc func(*packet.Packet) error
	size   int
}

// New returns a new packet factory
func New(lhost, rhost net.IP, lport, rport uint16, size int, fw func(*packet.Packet) error) (*PacketFactory, error) {
	if size > packet.MaxPayloadBytes {
		return nil, fmt.Errorf("max size is %d", packet.MaxPayloadBytes)
	}
	return &PacketFactory{
		lhost:  lhost,
		rhost:  rhost,
		lport:  lport,
		rport:  rport,
		fwFunc: fw,
		size:   size,
	}, nil
}

// DefaultPacketFactory returns a new packet factory with the maximum chunk size
func DefaultPacketFactory(lhost, rhost net.IP, lport, rport uint16, fw func(*packet.Packet) error) *PacketFactory {
	return &PacketFactory{
		lhost:  lhost,
		rhost:  rhost,
		lport:  lport,
		rport:  rport,
		fwFunc: fw,
		size:   packet.MaxPayloadBytes,
	}
}

// SendControlPacket crafts and sends a control packet to the network
func (pf *PacketFactory) SendControlPacket(syn, ack, fin, err bool) error {
	p, _ := packet.NewPacket(pf.lport, pf.rport, nil) // err checks for payload size (no payload)

	if syn {
		p.SetFlagSYN()
	}
	if ack {
		p.SetFlagACK()
	}
	if fin {
		p.SetFlagFIN()
	}
	if err {
		p.SetFlagERR()
	}
	p.SetSourceIPv4(pf.lhost)
	p.SetDestinationIPv4(pf.rhost)
	p.SetSum()

	if fwErr := pf.fwFunc(p); fwErr != nil {
		return fmt.Errorf("could not send control message SYN[%t] ACK[%t] FIN[%t] ERR[%t]: %s", syn, ack, fin, err, fwErr)
	}

	return nil
}

// PackAndForwardMessage chops a stream of bytes onto chunks of maximum size,
// wraps them in rdtp Packets and forwards them to the fwFunc
func (pf *PacketFactory) PackAndForwardMessage(msg []byte) (int, error) {
	var chunk []byte

	rem := msg
	txBytes := 0

	for len(rem) > 0 {
		if len(rem) >= pf.size {
			chunk, rem = rem[:pf.size], rem[pf.size:]
		} else {
			chunk, rem = rem, []byte{}
		}
		if err := pf.packetizeAndForwardChunk(chunk); err != nil {
			return txBytes, errors.Wrap(err, "could not packatize and forward chunk")
		}
		txBytes += len(chunk)
	}
	return txBytes, nil
}

func (pf *PacketFactory) packetizeAndForwardChunk(chunk []byte) error {
	pck, err := packet.NewPacket(pf.lport, pf.rport, chunk)
	if err != nil {
		return errors.Wrap(err, "error packetizing message")
	}
	pck.SetSourceIPv4(pf.lhost)
	pck.SetDestinationIPv4(pf.rhost)
	pck.SetSum() // set checksum here
	if err = pf.fwFunc(pck); err != nil {
		return errors.Wrap(err, "error forwarding packet")
	}
	return nil
}
