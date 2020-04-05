package atc

import (
	"sync"
	"time"

	"github.com/adrianosela/rdtp/packet"
)

var defaultAckWaitTime = time.Second * 1

// AirTrafficCtrl is the rdtp transmissions controller.
// It keeps track of packets transmitted but not acknowledged
// such that if the ack-wait timer times out, the packet will
// be retransmitted automatically.
type AirTrafficCtrl struct {
	sync.RWMutex // inherit read/write lock behavior

	ackWait time.Duration
	fwFunc  func(*packet.Packet) error

	inFlight map[uint32]*packet.Packet
}

// NewAirTrafficCtrl returns the default ATC
func NewAirTrafficCtrl(fwFunc func(*packet.Packet) error) *AirTrafficCtrl {
	return &AirTrafficCtrl{
		ackWait:  defaultAckWaitTime,
		fwFunc:   fwFunc,
		inFlight: make(map[uint32]*packet.Packet),
	}
}

// Send sends a packet while keeping track of it
func (atc *AirTrafficCtrl) Send(pck *packet.Packet) error {
	atc.Lock()
	defer atc.Unlock()

	atc.inFlight[pck.SeqNo] = pck

	return atc.fwFunc(pck)
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
