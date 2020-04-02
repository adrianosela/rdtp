package ctrl

import (
	"time"

	"github.com/adrianosela/rdtp/packet"
)

// Controller is the rdtp transmissions controller.
// It keeps track of packets transmitted but not acknowledged
// such that if the ack-wait timer times out, the packet will
// be retransmitted automatically.
type Controller interface {
	Send(*packet.Packet)
	Ack(uint16)
	SetAckWait(time.Duration)
}
