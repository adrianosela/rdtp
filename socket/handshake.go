package socket

import (
	"github.com/adrianosela/rdtp/handshake"
)

// Dial sends a SYN, waits for a SYN ACK, and sends an ACK
func (s *Socket) Dial() error {
	return handshake.InitiateConnection(s.inbound, s.packetizer.SendControlPacket)
}

// Accept sends a SYN ACK and waits for an ACK
func (s *Socket) Accept() error {
	return handshake.AcceptConnection(s.inbound, s.packetizer.SendControlPacket)
}

// finish manages the termination handshake
func (s *Socket) finish() error {
	select {
	case <-s.fin:
		return handshake.AcceptDisconnection(s.inbound, s.packetizer.SendControlPacket)
	default:
		return handshake.InitiateDisconnection(s.inbound, s.packetizer.SendControlPacket)
	}
}
