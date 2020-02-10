package multiplexer

import (
	"net"

	"github.com/adrianosela/rdtp/packet"
)

// Mux is the RDTP packet multiplexer
type Mux interface {
	Get(p uint16) (net.Conn, bool)
	Attach(p uint16, c net.Conn)
	Detach(p uint16)

	MultiplexPacket(p *packet.Packet) error
}
