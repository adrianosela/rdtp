package controller

import (
	"github.com/adrianosela/rdtp/packet"
	"github.com/adrianosela/rdtp/service/ports"
	"github.com/adrianosela/rdtp/socket"
)

// Controller represents the rdtp ports controller.
// It allocates and deallocates rdtp sockets and listeners.
type Controller interface {
	Put(sck *socket.Socket) error
	Evict(sckID string) error
	Deliver(p *packet.Packet) error
	AttachListener(l *ports.Listener) error
	DetachListener(port uint16) error
}
