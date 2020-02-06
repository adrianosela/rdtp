package controller

import (
	"fmt"
	"log"

	"github.com/adrianosela/rdtp"
)

// MultiplexPacket delivers a packet to the correct destination worker
func (ctrl *Controller) MultiplexPacket(p *rdtp.Packet) error {
	ctrl.RLock()
	defer ctrl.RUnlock()

	w, ok := ctrl.Ports[p.DstPort]
	if !ok {
		return fmt.Errorf("port %d is closed", p.DstPort)
	}

	w.rxChan <- p.Payload
	log.Printf("[MUX] %d --> %d", p.Length, p.DstPort)

	return nil
}
