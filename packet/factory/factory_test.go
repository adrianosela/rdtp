package factory

import (
	"fmt"
	"net"
	"testing"

	"github.com/adrianosela/rdtp/packet"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

var (
	testSrcIP = net.ParseIP("10.0.0.94")
	testDstIP = net.ParseIP("10.0.0.95")

	testMsg = []byte(
		"The Lord of the Rings is a film series of three epic fantasy adventure " +
			"films directed by Peter Jackson, based on the novel written by J. R. " +
			"R. Tolkien. The films are subtitled The Fellowship of the Ring (2001)")
)

func TestNewPacketFactoryOK(t *testing.T) {
	p, err := New(testSrcIP, testDstIP, 1234, 5678, packet.MaxPayloadBytes,
		func(x *packet.Packet) error {
			return nil
		})

	assert.NotNil(t, p)
	assert.Nil(t, err)
	assert.Equal(t, p.lport, uint16(1234))
	assert.Equal(t, p.rport, uint16(5678))
	assert.Equal(t, p.lhost, testSrcIP)
	assert.Equal(t, p.rhost, testDstIP)
	assert.Equal(t, p.size, packet.MaxPayloadBytes)
}

func TestNewPacketFactoryError(t *testing.T) {
	p, err := New(testSrcIP, testDstIP, 1234, 5678, packet.MaxPayloadBytes+1,
		func(x *packet.Packet) error {
			return nil
		})

	assert.Nil(t, p)
	assert.NotNil(t, err)
	assert.Equal(t, fmt.Errorf("max size is %d", packet.MaxPayloadBytes), err)
}

func TestDefaultPacketFactoryOK(t *testing.T) {
	p := DefaultPacketFactory(testSrcIP, testDstIP, 1234, 5678,
		func(x *packet.Packet) error {
			return nil
		})

	assert.NotNil(t, p)
	assert.Equal(t, p.lport, uint16(1234))
	assert.Equal(t, p.rport, uint16(5678))
	assert.Equal(t, p.lhost, testSrcIP)
	assert.Equal(t, p.rhost, testDstIP)
	assert.Equal(t, p.size, packet.MaxPayloadBytes)
}

func TestRunPacketFactoryFuncOK(t *testing.T) {
	funcRan := false

	p, err := New(testSrcIP, testDstIP, 1234, 5678, packet.MaxPayloadBytes,
		func(x *packet.Packet) error {
			funcRan = true
			return nil
		})

	assert.Nil(t, err)
	assert.False(t, funcRan)
	p.fwFunc(nil)
	assert.True(t, funcRan)
}

func TestSendControlPacketOK(t *testing.T) {
	var forwarded *packet.Packet

	pf := DefaultPacketFactory(testSrcIP, testDstIP, 1234, 5678,
		func(p *packet.Packet) error {
			forwarded = p
			return nil
		})
	assert.NotNil(t, pf)

	tests := []struct {
		syn bool
		ack bool
		fin bool
		err bool
	}{
		{syn: false, ack: false, fin: false, err: false},
		{syn: false, ack: false, fin: false, err: true},
		{syn: false, ack: false, fin: true, err: false},
		{syn: false, ack: false, fin: true, err: true},
		{syn: false, ack: true, fin: false, err: false},
		{syn: false, ack: true, fin: false, err: true},
		{syn: false, ack: true, fin: true, err: false},
		{syn: false, ack: true, fin: true, err: true},
		{syn: true, ack: false, fin: false, err: false},
		{syn: true, ack: false, fin: false, err: true},
		{syn: true, ack: false, fin: true, err: false},
		{syn: true, ack: false, fin: true, err: true},
		{syn: true, ack: true, fin: false, err: false},
		{syn: true, ack: true, fin: false, err: true},
		{syn: true, ack: true, fin: true, err: false},
		{syn: true, ack: true, fin: true, err: true},
	}

	for _, test := range tests {
		err := pf.SendControlPacket(test.syn, test.ack, test.fin, test.err)
		assert.Nil(t, err)
		assert.NotNil(t, forwarded)
		assert.Equal(t, forwarded.IsSYN(), test.syn)
		assert.Equal(t, forwarded.IsACK(), test.ack)
		assert.Equal(t, forwarded.IsFIN(), test.fin)
		assert.Equal(t, forwarded.IsERR(), test.err)
	}
}

func TestSendControlPacketError(t *testing.T) {
	mockError := errors.New("mock error")

	pf := DefaultPacketFactory(testSrcIP, testDstIP, 1234, 5678,
		func(p *packet.Packet) error {
			return mockError
		})
	assert.NotNil(t, pf)

	tests := []struct {
		syn bool
		ack bool
		fin bool
		err bool
	}{
		{syn: false, ack: false, fin: false, err: false},
		{syn: false, ack: false, fin: false, err: true},
		{syn: false, ack: false, fin: true, err: false},
		{syn: false, ack: false, fin: true, err: true},
		{syn: false, ack: true, fin: false, err: false},
		{syn: false, ack: true, fin: false, err: true},
		{syn: false, ack: true, fin: true, err: false},
		{syn: false, ack: true, fin: true, err: true},
		{syn: true, ack: false, fin: false, err: false},
		{syn: true, ack: false, fin: false, err: true},
		{syn: true, ack: false, fin: true, err: false},
		{syn: true, ack: false, fin: true, err: true},
		{syn: true, ack: true, fin: false, err: false},
		{syn: true, ack: true, fin: false, err: true},
		{syn: true, ack: true, fin: true, err: false},
		{syn: true, ack: true, fin: true, err: true},
	}

	for _, test := range tests {
		err := pf.SendControlPacket(test.syn, test.ack, test.fin, test.err)
		assert.NotNil(t, err)

		expectedError := fmt.Errorf(
			"could not send control message SYN[%t] ACK[%t] FIN[%t] ERR[%t]: %s",
			test.syn,
			test.ack,
			test.fin,
			test.err,
			mockError)

		assert.Equal(t, expectedError, err)
	}
}

func TestPacketizeAndForwardChunkOK(t *testing.T) {
	var rx []byte

	// subset of message
	size := len(testMsg) / 10
	chunk := testMsg[:size]

	p, err := New(testSrcIP, testDstIP, 1234, 5678, size,
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

	p, err := New(testSrcIP, testDstIP, 1234, 5678, 10, func(x *packet.Packet) error { return nil })
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

	p, err := New(testSrcIP, testDstIP, 1234, 5678, 3,
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

	p, err := New(testSrcIP, testDstIP, 1234, 5678, 5,
		func(x *packet.Packet) error {
			return mockError
		})

	assert.Nil(t, err)

	_, err = p.PackAndForwardMessage(testMsg)
	assert.NotNil(t, err)
	assert.Equal(t, errors.Wrap(mockError, "could not packatize and forward chunk: error forwarding packet").Error(), err.Error())
}
