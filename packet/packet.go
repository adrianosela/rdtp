package packet

import (
	"fmt"
)

const (
	// MaxPacketBytes is the maximum size of an RDTP packet incl. header
	MaxPacketBytes = 1500 // will chunk otherwise

	// HeaderByteSize is the byte size of an RDTP header
	HeaderByteSize = 9

	// MaxPayloadByteSize is the maximum size of a payload that a single RDTP
	// packet can carry
	MaxPayloadByteSize = MaxPacketBytes - HeaderByteSize

	// MaxPortNo is the maximum port number representable with 16 bits
	MaxPortNo = 65535
)

// Packet is an RDTP packet
type Packet struct {
	// connection identifyers
	SrcPort uint16
	DstPort uint16

	// processing and integrity
	Length   uint16
	Checksum uint16

	Flags uint8 // {SYN, FIN, ACK, ERR, XXXX, XXXX, XXXX, XXXX}

	Payload []byte
}

// NewPacket populates an RDTP packet onto a serializable state representation
func NewPacket(src, dst uint16, payload []byte) (*Packet, error) {
	if len(payload) > MaxPayloadByteSize {
		return nil, fmt.Errorf(
			"Invalid RDTP payload. Payload length %d more than %d bytes",
			len(payload),
			MaxPayloadByteSize,
		)
	}
	p := &Packet{
		SrcPort: src,
		DstPort: dst,
		Length:  uint16(len(payload)),
		Flags:   uint8(0),
		Payload: payload,
	}
	p.Checksum = p.computeChecksum()
	return p, nil
}
