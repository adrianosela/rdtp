package rdtp

import (
	"fmt"
)

// Packet is an RDTP packet
type Packet struct {
	SrcPort  uint16
	DstPort  uint16
	Length   uint16
	Checksum uint16
	Payload  []byte
}

const (
	// MaxPacketBytes is the maximum size of an RDTP packet incl. header
	MaxPacketBytes = 65515 // 65535 - IPv4 Header (20 bytes)
	// HeaderByteSize is the byte size of an RDTP header
	HeaderByteSize = 8
	// MaxPayloadByteSize is the maximum size of a payload that a single RDTP
	// packet can carry
	MaxPayloadByteSize = MaxPacketBytes - HeaderByteSize
)

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
		Length:  uint16(HeaderByteSize + len(payload)),
		Payload: payload,
	}
	p.Checksum = p.computeChecksum()
	return p, nil
}
