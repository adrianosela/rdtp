package rtx

import (
	"time"

	"github.com/adrianosela/rdtp/packet"
)

// Controller is the RDTP retransmissions controller.
// It keeps track of packets transmitted but not acknowledged
// such that if the ack-wait timer times out, the packet will
// be retransmitted.
type Controller interface {
	// send a packet
	Send(*packet.Packet)

	// ack a packet
	Ack(uint16)

	// set ack-wait timer time
	SetAckWait(time.Duration)
}
