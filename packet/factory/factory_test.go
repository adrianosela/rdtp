package factory

import (
	"testing"

	"github.com/adrianosela/rdtp/packet"
	"github.com/stretchr/testify/assert"
)

var testMsg = []byte(
	"The Lord of the Rings is a film series of three epic fantasy adventure " +
		"films directed by Peter Jackson, based on the novel written by J. R. " +
		"R. Tolkien. The films are subtitled The Fellowship of the Ring (2001)")

func TestNewPacketizer(t *testing.T) {
	p, err := New(1234, 5678,
		func(x *packet.Packet) error {
			return nil
		}, packet.MaxPayloadBytes)

	assert.Nil(t, err)
	assert.Equal(t, p.srcPort, uint16(1234))
	assert.Equal(t, p.dstPort, uint16(5678))
	assert.Equal(t, p.size, packet.MaxPayloadBytes)
}

func TestRunPacketizerFunc(t *testing.T) {
	funcRan := false

	p, err := New(1234, 5678,
		func(x *packet.Packet) error {
			funcRan = true
			return nil
		}, packet.MaxPayloadBytes)

	assert.Nil(t, err)
	assert.False(t, funcRan)
	p.fwFunc(nil)
	assert.True(t, funcRan)
}

func TestPacketizeAndForwardChunk(t *testing.T) {
	var rx []byte

	// subset of message
	size := len(testMsg) / 10
	chunk := testMsg[:size]

	p, err := New(1234, 5678,
		func(x *packet.Packet) error {
			rx = append(rx, x.Payload...)
			return nil
		}, size)
	assert.Nil(t, err)

	err = p.packetizeAndForwardChunk(chunk)
	assert.Nil(t, err)

	// check chunk sent and received match
	assert.Equal(t, chunk, rx)
}

func TestSend(t *testing.T) {
	var rx []byte

	p, err := New(1234, 5678,
		func(x *packet.Packet) error {
			rx = append(rx, x.Payload...)
			return nil
		}, 3)

	assert.Nil(t, err)

	n, err := p.Send(testMsg)

	// check sent the whole message
	assert.Nil(t, err)
	assert.Equal(t, len(testMsg), n)

	// check message sent and received match
	assert.Equal(t, testMsg, rx)
}
