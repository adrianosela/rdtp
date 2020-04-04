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
func New(src, dst uint16, fw func(*packet.Packet) error, size int) (*PacketFactory, error) {
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

// Send chops a stream of bytes onto chunks of maximum size
// and forwards them as soon as they are ready
func (pf *PacketFactory) Send(msg []byte) (int, error) {
	var chunk []byte

	rem := msg
	txBytes := 0

	for len(rem) >= pf.size {
		chunk, rem = rem[:pf.size], rem[pf.size:]
		if err := pf.packetizeAndForwardChunk(chunk); err != nil {
			return txBytes, errors.Wrap(err, "could not packatize and forward chunk")
		}
		txBytes += len(chunk)
	}

	if len(rem) > 0 {
		if err := pf.packetizeAndForwardChunk(rem); err != nil {
			return txBytes, errors.Wrap(err, "could not packatize and forward chunk")
		}
		txBytes += len(rem)
	}

	return txBytes, nil
}

func (pf *PacketFactory) packetizeAndForwardChunk(chunk []byte) error {
	pck, err := packet.NewPacket(pf.srcPort, pf.dstPort, chunk)
	if err != nil {
		return errors.Wrap(err, "error packetizing message")
	}
	if err = pf.fwFunc(pck); err != nil {
		return errors.Wrap(err, "error forwarding packet")
	}
	return nil
}
