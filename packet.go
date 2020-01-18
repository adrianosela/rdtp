package rdtp

import (
	"encoding/binary"
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
)

// NewPacket populates an RDTP packet onto a serializable state representation
func NewPacket(src, dst uint16, payload []byte) (*Packet, error) {
	if len(payload) > MaxPacketBytes-HeaderByteSize {
		return nil, fmt.Errorf(
			"Invalid RDTP payload. Payload length %d more than %d bytes",
			len(payload),
			MaxPacketBytes-HeaderByteSize,
		)
	}
	p := &Packet{
		SrcPort: src,
		DstPort: dst,
		Length:  uint16(HeaderByteSize + len(payload)),
		Payload: payload,
	}
	p.Checksum = computeChecksum(p)
	return p, nil
}

// Serialize byte-encodes an RDTP packet ready to be encapsulated
// in a network layer protocol packet (i.e. IP datagram)
func (p *Packet) Serialize() []byte {
	b := make([]byte, HeaderByteSize)
	binary.BigEndian.PutUint16(b[0:2], p.SrcPort)
	binary.BigEndian.PutUint16(b[2:4], p.DstPort)
	binary.BigEndian.PutUint16(b[4:6], p.Length)
	binary.BigEndian.PutUint16(b[6:8], p.Checksum)
	return append(b, p.Payload...)
}

// Deserialize byte decodes an RDTP packet
func Deserialize(data []byte) (*Packet, error) {
	if len(data) < HeaderByteSize {
		return nil, fmt.Errorf(
			"Invalid RDTP header. Packet length %d less than %d bytes",
			len(data),
			HeaderByteSize)
	}
	p := &Packet{
		SrcPort:  binary.BigEndian.Uint16(data[0:2]),
		DstPort:  binary.BigEndian.Uint16(data[2:4]),
		Length:   binary.BigEndian.Uint16(data[4:6]),
		Checksum: binary.BigEndian.Uint16(data[6:8]),
		Payload:  data[8:],
	}
	return p, nil
}
