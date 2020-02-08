package controller

import (
	"fmt"

	"github.com/adrianosela/rdtp/packet"
)

// MultiplexPacket delivers a packet to the correct destination worker
func (ctrl *Controller) MultiplexPacket(p *packet.Packet) error {
	ctrl.RLock()
	defer ctrl.RUnlock()

	worker, ok := ctrl.Ports[p.DstPort]
	if !ok {
		return fmt.Errorf("port %d is closed", p.DstPort)
	}

	worker.rxChan <- p.Payload

	return nil
}
