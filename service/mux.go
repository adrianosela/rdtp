package service

import (
	"log"

	"github.com/adrianosela/rdtp/packet"
)

// MultiplexPacket delivers a packet to the correct destination worker
func (s *Service) MultiplexPacket(p *packet.Packet) error {
	// TODO
	log.Printf("[MUX] RX %d ==> %d", p.Length, p.DstPort)
	return nil
}
