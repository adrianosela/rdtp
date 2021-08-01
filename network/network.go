package network

import (
	"github.com/adrianosela/rdtp/packet"
)

// Network represents an unreliable channel for sending and receiving rdtp packets
type Network interface {
	Send(p *packet.Packet) error
	StartReceiver(fn func(p *packet.Packet) error)
}
