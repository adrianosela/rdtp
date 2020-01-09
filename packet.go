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

// NewPacket populates an RDTP packet onto a serializable state representation
func NewPacket(src, dst uint16, payload []byte) (*Packet, error) {
	if len(payload) > MaxPacketBytes-MinHeaderBytes {
		return nil, fmt.Errorf(
			"Invalid RDTP payload. Payload length %d more than %d bytes",
			len(payload),
			MaxPacketBytes-MinHeaderBytes,
		)
	}
	p := &Packet{
		SrcPort:  src,
		DstPort:  dst,
		Length:   uint16(MinHeaderBytes + len(payload)),
		Checksum: uint16(0), // zero out initially
		Payload:  payload,
	}
	p.Checksum = computeChecksum(p)
	return p, nil
}

// Serialize byte-encodes an RDTP packet ready to be encapsulated
// in a network layer protocol packet (i.e. IP datagram)
func (p *Packet) Serialize() ([]byte, error) {
	// TODO
	return nil, nil
}

// Deserialize byte decodes an RDTP packet
func Deserialize(data []byte) (*Packet, error) {
	if len(data) < MinHeaderBytes {
		return nil, fmt.Errorf(
			"Invalid RDTP header. Packet length %d less than %d bytes",
			len(data),
			MinHeaderBytes)
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
