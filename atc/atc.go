package atc

import (
	"sync"
	"time"

	"github.com/adrianosela/rdtp/packet"
)

// AirTrafficCtrl is the rdtp transmissions controller.
// It keeps track of packets transmitted but not acknowledged
// such that if the ack-wait timer times out, the packet will
// be retransmitted automatically.
type AirTrafficCtrl struct {
	sync.RWMutex // inherit read/write lock behavior

	inFlight map[uint32]*packet.Packet
	ackWait  time.Duration
}

// Ack acknowledges a sent packet
func (atc *AirTrafficCtrl) Ack(num uint32) {
	atc.Lock()
	defer atc.Unlock()

	delete(atc.inFlight, num)
}

// SetAckWait sets the Ack wait timer time on the controller
func (atc *AirTrafficCtrl) SetAckWait(t time.Duration) {
	atc.Lock()
	defer atc.Unlock()

	atc.ackWait = t
}
