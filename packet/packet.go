package packet

import (
	"fmt"
	"net"
)

const (
	// MaxPacketBytes is the maximum size of an RDTP packet incl. header
	MaxPacketBytes = 1500 // will chunk otherwise

	// HeaderByteSize is the byte size of an RDTP header
	HeaderByteSize = 17

	// MaxPayloadBytes is the maximum size of a payload that
	// a single RDTP packet can carry
	MaxPayloadBytes = MaxPacketBytes - HeaderByteSize
)

// Packet is an RDTP packet
type Packet struct {
	// connection identifyers
	SrcPort uint16
	DstPort uint16

	// processing and integrity
	Length   uint16
	Checksum uint16

	// reliability
	SeqNo uint32
	AckNo uint32

	// control
	Flags uint8 // {SYN, FIN, ACK, ERR, XXXX, XXXX, XXXX, XXXX}

	// data
	Payload []byte

	// the fields below dont make up the
	// packet that goes over the wire.
	// they are used to communicate
	// network layer metadata
	srcIP net.IP
	dstIP net.IP
}

// NewPacket populates an RDTP packet onto a serializable state representation
func NewPacket(src, dst uint16, payload []byte) (*Packet, error) {
	if len(payload) > MaxPayloadBytes {
		return nil, fmt.Errorf(
			"invalid rdtp payload - payload length %d more than %d bytes",
			len(payload),
			MaxPayloadBytes,
		)
	}
	p := &Packet{
		SrcPort: src,
		DstPort: dst,
		Length:  uint16(len(payload)),
		Payload: payload,
	}
	return p, nil
}
