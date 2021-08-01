package factory

import (
	"fmt"

	"github.com/adrianosela/rdtp/packet"
	"github.com/pkg/errors"
)

// PacketFactory formats streams of arbitrary data onto
// chunks of a given max size, builds an rdtp packet
// with the data, and forwards it
type PacketFactory struct {
	srcPort uint16
	dstPort uint16
	fwFunc  func(*packet.Packet) error
	size    int
}

// New returns a new packet factory
func New(src, dst uint16, size int, fw func(*packet.Packet) error) (*PacketFactory, error) {
	if size > packet.MaxPayloadBytes {
		return nil, fmt.Errorf("max size is %d", packet.MaxPayloadBytes)
	}
	return &PacketFactory{
		srcPort: src,
		dstPort: dst,
		fwFunc:  fw,
		size:    size,
	}, nil
}

// DefaultPacketFactory returns a new packet factory with the maximum chunk size
func DefaultPacketFactory(src, dst uint16, fw func(*packet.Packet) error) *PacketFactory {
	return &PacketFactory{
		srcPort: src,
		dstPort: dst,
		fwFunc:  fw,
		size:    packet.MaxPayloadBytes,
	}
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
	pck, err := packet.NewPacket(pf.srcPort, pf.dstPort, chunk)
	if err != nil {
		return errors.Wrap(err, "error packetizing message")
	}
	pck.SetSum() // set checksum here
	if err = pf.fwFunc(pck); err != nil {
		return errors.Wrap(err, "error forwarding packet")
	}
	return nil
}
