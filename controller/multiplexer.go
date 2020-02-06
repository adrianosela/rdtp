package controller

import (
	"log"

	"github.com/adrianosela/rdtp"
)

// MultiplexPacket delivers a packet to the correct destination worker
func (ctrl *Controller) MultiplexPacket(p *rdtp.Packet) error {
	// TODO
	log.Printf("[MUX] %d --> %d", p.Length, p.DstPort)
	return nil
}
