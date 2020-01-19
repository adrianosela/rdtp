package rdtp

import (
	"encoding/binary"
	"fmt"
)

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
		Checksum: binary.BigEndian.Uint16(data[6:HeaderByteSize]),
		Payload:  data[HeaderByteSize:],
	}
	// safely clean up payload length
	if p.Length <= uint16(len(p.Payload)) {
		p.Payload = p.Payload[:p.Length]
	} else {
		return nil, fmt.Errorf(
			"Invalid RDTP header. 'Length' field (%d) longer than data (%d)",
			p.Length,
			len(data)-HeaderByteSize)
	}
	return p, nil
}
