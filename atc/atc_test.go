package atc

import (
	"fmt"
	"testing"

	"github.com/adrianosela/rdtp/packet"
	"github.com/stretchr/testify/assert"
)

var (
	mockPacket, _ = packet.NewPacket(uint16(10), uint16(20), []byte("hello"))
)

func TestNewAirTrafficCtrl(t *testing.T) {
	atc := NewAirTrafficCtrl(func(p *packet.Packet) error { return nil })
	assert.NotNil(t, atc)
	assert.Equal(t, defaultAckWaitTime, atc.ackWait)
	assert.Equal(t, make(map[uint32]*packet.Packet), atc.inFlight)
}

func TestSendNoError(t *testing.T) {
	var sentPacket *packet.Packet

	mockPacket.SetSeqNo(100)

	atc := NewAirTrafficCtrl(func(p *packet.Packet) error {
		sentPacket = p
		return nil
	})
	assert.NotNil(t, atc)

	assert.Nil(t, sentPacket)
	err := atc.Send(mockPacket)
	assert.Nil(t, err)
	assert.NotNil(t, sentPacket)
	assert.Equal(t, mockPacket, sentPacket)

	assert.Contains(t, atc.inFlight, mockPacket.SeqNo)
}

func TestSendWithError(t *testing.T) {
	var sentPacket *packet.Packet

	mockError := fmt.Errorf("mock error")

	atc := NewAirTrafficCtrl(func(p *packet.Packet) error {
		return mockError
	})
	assert.NotNil(t, atc)

	assert.Nil(t, sentPacket)
	err := atc.Send(mockPacket)
	assert.NotNil(t, err)
	assert.Equal(t, err, fmt.Errorf("could not send packet: %s", mockError))
	assert.Nil(t, sentPacket)
}

func TestAckSentPacket(t *testing.T) {
	var sentPacket *packet.Packet

	atc := NewAirTrafficCtrl(func(p *packet.Packet) error {
		sentPacket = p
		return nil
	})
	assert.NotNil(t, atc)

	assert.Nil(t, sentPacket)
	err := atc.Send(mockPacket)
	assert.Nil(t, err)
	assert.NotNil(t, sentPacket)
	assert.Equal(t, mockPacket, sentPacket)

	assert.Contains(t, atc.inFlight, mockPacket.SeqNo)
	atc.Ack(mockPacket.SeqNo)
	assert.NotContains(t, atc.inFlight, mockPacket.SeqNo)
}
