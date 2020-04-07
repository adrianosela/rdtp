package packet

import (
	"fmt"

	"github.com/google/gopacket/layers"
)

const (
	// MaxPacketBytes is the maximum size of an RDTP packet incl. header
	MaxPacketBytes = 1500 // will chunk otherwise

	// HeaderByteSize is the byte size of an RDTP header
	HeaderByteSize = 17

	// MaxPayloadBytes is the maximum size of a payload that a single RDTP
	// packet can carry
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

	Flags uint8 // {SYN, FIN, ACK, ERR, XXXX, XXXX, XXXX, XXXX}

	Payload []byte

	// used exclusively on inbound packets to
	// communicate important network layer
	// details. e.g. ip addresses for
	// identifying recieving socket
	ipv4 *layers.IPv4
}

// NewPacket populates an RDTP packet onto a serializable state representation
func NewPacket(src, dst uint16, payload []byte) (*Packet, error) {
	if len(payload) > MaxPayloadBytes {
		return nil, fmt.Errorf(
			"Invalid RDTP payload. Payload length %d more than %d bytes",
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
