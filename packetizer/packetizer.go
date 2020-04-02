package packetizer

import (
	"fmt"

	"github.com/adrianosela/rdtp/packet"
	"github.com/pkg/errors"
)

// Packetizer formats streams of arbitrary data onto
// chunks of a given max size, builds an rdtp packet
// with the data, and forwards it
type Packetizer struct {
	srcPort uint16
	dstPort uint16
	fwFunc  func(*packet.Packet) error
	size    int
}

// New returns a new packetizer
func New(src, dst uint16, fw func(*packet.Packet) error, size int) (*Packetizer, error) {
	if size > packet.MaxPayloadBytes {
		return nil, fmt.Errorf("max size is %d", packet.MaxPayloadBytes)
	}
	return &Packetizer{
		srcPort: src,
		dstPort: dst,
		fwFunc:  fw,
		size:    size,
	}, nil
}

// Send chops a stream of bytes onto chunks of maximum size
// and forwards them as soon as they are ready
func (p *Packetizer) Send(msg []byte) (int, error) {
	var chunk []byte

	rem := msg
	txBytes := 0

	for len(rem) >= p.size {
		chunk, rem = rem[:p.size], rem[p.size:]
		if err := p.packetizeAndForwardChunk(chunk); err != nil {
			return txBytes, errors.Wrap(err, "could not packatize and forward chunk")
		}
		txBytes += len(chunk)
	}

	if len(rem) > 0 {
		if err := p.packetizeAndForwardChunk(rem); err != nil {
			return txBytes, errors.Wrap(err, "could not packatize and forward chunk")
		}
		txBytes += len(rem)
	}

	return txBytes, nil
}

func (p *Packetizer) packetizeAndForwardChunk(chunk []byte) error {
	pck, err := packet.NewPacket(p.srcPort, p.dstPort, chunk)
	if err != nil {
		return errors.Wrap(err, "error packetizing message")
	}
	if err = p.fwFunc(pck); err != nil {
		return errors.Wrap(err, "error forwarding packet")
	}
	return nil
}
