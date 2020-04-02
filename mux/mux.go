package mux

import (
	"net"

	"github.com/adrianosela/rdtp/packet"
)

// Mux is the RDTP packet multiplexer
type Mux interface {
	Get(p uint16) (net.Conn, error)
	Attach(p uint16, c net.Conn) error
	AttachAny(c net.Conn) (uint16, error)
	Detach(p uint16)

	MultiplexPacket(p *packet.Packet) error
}
