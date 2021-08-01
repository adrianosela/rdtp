package factory

import (
	"fmt"
	"testing"

	"github.com/adrianosela/rdtp/packet"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

var testMsg = []byte(
	"The Lord of the Rings is a film series of three epic fantasy adventure " +
		"films directed by Peter Jackson, based on the novel written by J. R. " +
		"R. Tolkien. The films are subtitled The Fellowship of the Ring (2001)")

func TestNewPacketFactoryOK(t *testing.T) {
	p, err := New(1234, 5678, packet.MaxPayloadBytes,
		func(x *packet.Packet) error {
			return nil
		})

	assert.NotNil(t, p)
	assert.Nil(t, err)
	assert.Equal(t, p.srcPort, uint16(1234))
	assert.Equal(t, p.dstPort, uint16(5678))
	assert.Equal(t, p.size, packet.MaxPayloadBytes)
}

func TestNewPacketFactoryError(t *testing.T) {
	p, err := New(1234, 5678, packet.MaxPayloadBytes+1,
		func(x *packet.Packet) error {
			return nil
		})

	assert.Nil(t, p)
	assert.NotNil(t, err)
	assert.Equal(t, fmt.Errorf("max size is %d", packet.MaxPayloadBytes), err)
}

func TestDefaultPacketFactoryOK(t *testing.T) {
	p := DefaultPacketFactory(1234, 5678,
		func(x *packet.Packet) error {
			return nil
		})

	assert.NotNil(t, p)
	assert.Equal(t, p.srcPort, uint16(1234))
	assert.Equal(t, p.dstPort, uint16(5678))
	assert.Equal(t, p.size, packet.MaxPayloadBytes)
}

func TestRunPacketFactoryFuncOK(t *testing.T) {
	funcRan := false

	p, err := New(1234, 5678, packet.MaxPayloadBytes,
		func(x *packet.Packet) error {
			funcRan = true
			return nil
		})

	assert.Nil(t, err)
	assert.False(t, funcRan)
	p.fwFunc(nil)
	assert.True(t, funcRan)
}

func TestPacketizeAndForwardChunkOK(t *testing.T) {
	var rx []byte

	// subset of message
	size := len(testMsg) / 10
	chunk := testMsg[:size]

	p, err := New(1234, 5678, size,
		func(x *packet.Packet) error {
			rx = append(rx, x.Payload...)
			return nil
		})
	assert.Nil(t, err)

	err = p.packetizeAndForwardChunk(chunk)
	assert.Nil(t, err)

	// check chunk sent and received match
	assert.Equal(t, chunk, rx)
}

func TestPacketizeAndForwardChunkError(t *testing.T) {
	badChunkLength := packet.MaxPayloadBytes + 1

	chunk := make([]byte, badChunkLength)

	p, err := New(1234, 5678, 10, func(x *packet.Packet) error { return nil })
	assert.Nil(t, err)

	err = p.packetizeAndForwardChunk(chunk)
	assert.NotNil(t, err)
	assert.Equal(t,
		fmt.Errorf(
			"error packetizing message: invalid rdtp payload - payload length %d more than %d bytes",
			badChunkLength,
			packet.MaxPayloadBytes).Error(),
		err.Error())
}

func TestPackAndForwardMessageOK(t *testing.T) {
	var rx []byte

	p, err := New(1234, 5678, 3,
		func(x *packet.Packet) error {
			rx = append(rx, x.Payload...)
			return nil
		})

	assert.Nil(t, err)

	n, err := p.PackAndForwardMessage(testMsg)

	// check sent the whole message
	assert.Nil(t, err)
	assert.Equal(t, len(testMsg), n)

	// check message sent and received match
	assert.Equal(t, testMsg, rx)
}

func TestPackAndForwardMessageError(t *testing.T) {
	mockError := errors.New("mock error")

	p, err := New(1234, 5678, 5,
		func(x *packet.Packet) error {
			return mockError
		})

	assert.Nil(t, err)

	_, err = p.PackAndForwardMessage(testMsg)
	assert.NotNil(t, err)
	assert.Equal(t, errors.Wrap(mockError, "could not packatize and forward chunk: error forwarding packet").Error(), err.Error())
}
