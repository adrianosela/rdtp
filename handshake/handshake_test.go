package handshake

import (
	"errors"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/adrianosela/rdtp/packet"
	"github.com/stretchr/testify/assert"
)

var (
	errMock = errors.New("mock error")
	msgMock = "mock message"

	flagCombinations = []struct {
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
)

func TestInitiateConnectionOK(t *testing.T) {
	local := make(chan *packet.Packet)
	remote := make(chan *packet.Packet)

	// mock remote operation
	go func() {
		// wait for SYN
		p := <-remote
		assert.True(t, p.IsSYN())
		assert.False(t, p.IsACK())
		assert.False(t, p.IsFIN())
		assert.False(t, p.IsERR())

		// send SYN ACK
		local <- mockControlPacket(true, true, false, false)

		// wait for ACK
		p = <-remote
		assert.False(t, p.IsSYN())
		assert.True(t, p.IsACK())
		assert.False(t, p.IsFIN())
		assert.False(t, p.IsERR())
	}()

	err := InitiateConnection(local, func(syn, ack, fin, err bool) error {
		remote <- mockControlPacket(syn, ack, fin, err)
		return nil
	})
	assert.Nil(t, err)

	// let mock remote go routine complete
	time.Sleep(time.Millisecond * 1)
}

func TestInitiateConnectionSendSynError(t *testing.T) {
	err := InitiateConnection(make(chan *packet.Packet), func(syn, ack, fin, err bool) error {
		return errMock
	})
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), fmt.Sprintf("connect handshake failed when sending SYN: %s", errMock))
}

func TestInitiateConnectionWaitForSynAckError(t *testing.T) {
	local := make(chan *packet.Packet)
	remote := make(chan *packet.Packet)

	// mock remote operation
	go func() {
		// wait for SYN
		p := <-remote
		assert.True(t, p.IsSYN())
		assert.False(t, p.IsACK())
		assert.False(t, p.IsFIN())
		assert.False(t, p.IsERR())

		// don't send SYN ACK (let it time out)
	}()

	err := InitiateConnection(local, func(syn, ack, fin, err bool) error {
		remote <- mockControlPacket(syn, ack, fin, err)
		return nil
	})
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "connect handshake failed when waiting for SYN ACK: operation timed out")
}

func TestInitiateConnectionSendAckError(t *testing.T) {
	local := make(chan *packet.Packet)
	remote := make(chan *packet.Packet)

	// mock remote operation
	go func() {
		// wait for SYN
		p := <-remote
		assert.True(t, p.IsSYN())
		assert.False(t, p.IsACK())
		assert.False(t, p.IsFIN())
		assert.False(t, p.IsERR())

		// send SYN ACK
		local <- mockControlPacket(true, true, false, false)

		// dont wait for ACK (local fails)
	}()

	sendInvocations := 0

	err := InitiateConnection(local, func(syn, ack, fin, err bool) error {
		if sendInvocations > 0 {
			return errMock
		}

		remote <- mockControlPacket(syn, ack, fin, err)
		sendInvocations++
		return nil
	})

	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), fmt.Sprintf("connect handshake failed when sending ACK: %s", errMock))
}

func TestAcceptConnectionOK(t *testing.T) {
	local := make(chan *packet.Packet)
	remote := make(chan *packet.Packet)

	// mock remote operation
	go func() {
		// AcceptConnection is only invoked upon receiving
		// a SYN packet so no need to send one here

		// wait for SYN ACK
		p := <-remote
		assert.True(t, p.IsSYN())
		assert.True(t, p.IsACK())
		assert.False(t, p.IsFIN())
		assert.False(t, p.IsERR())

		// send ACK
		local <- mockControlPacket(false, true, false, false)
	}()

	err := AcceptConnection(local, func(syn, ack, fin, err bool) error {
		remote <- mockControlPacket(syn, ack, fin, err)
		return nil
	})
	assert.Nil(t, err)
}

func TestAcceptConnectionSendSynAckError(t *testing.T) {
	err := AcceptConnection(make(chan *packet.Packet), func(syn, ack, fin, err bool) error {
		return errMock
	})
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), fmt.Sprintf("connect handshake failed when sending SYN ACK: %s", errMock))
}

func TestAcceptConnectionWaitForAckError(t *testing.T) {
	remote := make(chan *packet.Packet)

	// mock remote operation
	go func() {
		// wait for SYN ACK
		p := <-remote
		assert.True(t, p.IsSYN())
		assert.True(t, p.IsACK())
		assert.False(t, p.IsFIN())
		assert.False(t, p.IsERR())

		// don't send ACK (let it time out)
	}()

	err := AcceptConnection(make(chan *packet.Packet), func(syn, ack, fin, err bool) error {
		remote <- mockControlPacket(syn, ack, fin, err)
		return nil
	})
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "connect handshake failed when waiting for ACK: operation timed out")
}

func TestInitiateDisconnectionOK(t *testing.T) {
	local := make(chan *packet.Packet)
	remote := make(chan *packet.Packet)

	// mock remote operation
	go func() {
		// wait for FIN
		p := <-remote
		assert.False(t, p.IsSYN())
		assert.False(t, p.IsACK())
		assert.True(t, p.IsFIN())
		assert.False(t, p.IsERR())

		// send FIN ACK
		local <- mockControlPacket(false, true, true, false)

		// wait for ACK
		p = <-remote
		assert.False(t, p.IsSYN())
		assert.True(t, p.IsACK())
		assert.False(t, p.IsFIN())
		assert.False(t, p.IsERR())
	}()

	err := InitiateDisconnection(local, func(syn, ack, fin, err bool) error {
		remote <- mockControlPacket(syn, ack, fin, err)
		return nil
	})
	assert.Nil(t, err)

	// let mock remote go routine complete
	time.Sleep(time.Millisecond * 1)
}

func TestInitiateDisconnectionSendFinError(t *testing.T) {
	err := InitiateDisconnection(make(chan *packet.Packet), func(syn, ack, fin, err bool) error {
		return errMock
	})
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), fmt.Sprintf("finish handshake failed when sending FIN: %s", errMock))
}

func TestInitiateDisconnectionWaitForFinAckError(t *testing.T) {
	local := make(chan *packet.Packet)
	remote := make(chan *packet.Packet)

	// mock remote operation
	go func() {
		// wait for FIN
		p := <-remote
		assert.False(t, p.IsSYN())
		assert.False(t, p.IsACK())
		assert.True(t, p.IsFIN())
		assert.False(t, p.IsERR())

		// don't send FIN ACK (let it time out)
	}()

	err := InitiateDisconnection(local, func(syn, ack, fin, err bool) error {
		remote <- mockControlPacket(syn, ack, fin, err)
		return nil
	})
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "finish handshake failed when waiting for FIN ACK: operation timed out")
}

func TestInitiateDisconnectionSendAckError(t *testing.T) {
	local := make(chan *packet.Packet)
	remote := make(chan *packet.Packet)

	// mock remote operation
	go func() {
		// wait for FIN
		p := <-remote
		assert.False(t, p.IsSYN())
		assert.False(t, p.IsACK())
		assert.True(t, p.IsFIN())
		assert.False(t, p.IsERR())

		// send FIN ACK
		local <- mockControlPacket(false, true, true, false)

		// dont wait for ACK (local fails)
	}()

	sendInvocations := 0

	err := InitiateDisconnection(local, func(syn, ack, fin, err bool) error {
		if sendInvocations > 0 {
			return errMock
		}

		remote <- mockControlPacket(syn, ack, fin, err)
		sendInvocations++
		return nil
	})

	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), fmt.Sprintf("finish handshake failed when sending ACK: %s", errMock))
}

func TestAcceptDisconnectionOK(t *testing.T) {
	local := make(chan *packet.Packet)
	remote := make(chan *packet.Packet)

	// mock remote operation
	go func() {
		// AcceptDisonnection is only invoked upon receiving
		// a fin packet so no need to send one here

		// wait for fin ACK
		p := <-remote
		assert.False(t, p.IsSYN())
		assert.True(t, p.IsACK())
		assert.True(t, p.IsFIN())
		assert.False(t, p.IsERR())

		// send ACK
		local <- mockControlPacket(false, true, false, false)
	}()

	err := AcceptDisconnection(local, func(syn, ack, fin, err bool) error {
		remote <- mockControlPacket(syn, ack, fin, err)
		return nil
	})
	assert.Nil(t, err)
}

func TestAcceptDisconnectionSendFinAckError(t *testing.T) {
	err := AcceptDisconnection(make(chan *packet.Packet), func(syn, ack, fin, err bool) error {
		return errMock
	})
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), fmt.Sprintf("finish handshake failed when sending FIN ACK: %s", errMock))
}

func TestAcceptDisconnectionWaitForAckError(t *testing.T) {
	remote := make(chan *packet.Packet)

	// mock remote operation
	go func() {
		// wait for FIN ACK
		p := <-remote
		assert.False(t, p.IsSYN())
		assert.True(t, p.IsACK())
		assert.True(t, p.IsFIN())
		assert.False(t, p.IsERR())

		// don't send ACK (let it time out)
	}()

	err := AcceptDisconnection(make(chan *packet.Packet), func(syn, ack, fin, err bool) error {
		remote <- mockControlPacket(syn, ack, fin, err)
		return nil
	})
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "finish handshake failed when waiting for ACK: operation timed out")
}

func TestReceiveControlPacketOK(t *testing.T) {
	for _, comb := range flagCombinations {
		recvChan := make(chan *packet.Packet)
		go func() {
			recvChan <- mockControlPacket(comb.syn, comb.ack, comb.fin, comb.err)
		}()
		err := receiveControlPacket(
			recvChan,
			comb.syn, comb.ack, comb.fin, comb.err,
			time.Millisecond*1 /* no network inbetween -- use short timeout */)
		assert.Nil(t, err)
	}
}

func TestReceiveControlPacketErrorTimeout(t *testing.T) {
	for _, comb := range flagCombinations {
		recvChan := make(chan *packet.Packet)
		err := receiveControlPacket(
			recvChan,
			comb.syn, comb.ack, comb.fin, comb.err,
			time.Nanosecond*1 /* no network inbetween -- use short timeout */)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), "operation timed out")
	}
}

func TestReceiveControlPacketErrorUnexpectedPacket(t *testing.T) {
	recvChan := make(chan *packet.Packet)

	for _, expect := range flagCombinations {
		for _, get := range flagCombinations {
			go func() {
				recvChan <- mockControlPacket(get.syn, get.ack, get.fin, get.err)
			}()

			err := receiveControlPacket(recvChan, expect.syn, expect.ack, expect.fin, expect.err, time.Millisecond*1)

			if expect.syn != get.syn || expect.ack != get.ack || expect.fin != get.fin || expect.err != get.err {
				assert.NotNil(t, err)
				assert.Equal(t, err.Error(), fmt.Sprintf(
					"expected packet with flags %s but got %s",
					fmt.Sprintf(flagFmt, expect.syn, expect.ack, expect.fin, expect.err),
					fmt.Sprintf(flagFmt, get.syn, get.ack, get.fin, get.err)))
			} else {
				assert.Nil(t, err)
			}
		}
	}
}

type assertingWriter struct {
	t       *testing.T
	onWrite func()
	msg     string
}

func (w assertingWriter) Write(data []byte) (int, error) {
	w.onWrite()
	// assert message written after the dateTime (20 bytes), ignore last byte (\n)
	assert.Equal(w.t, w.msg, string(data[20:len(data)-1]))
	return len(data), nil
}

func TestConditionallyLogTrue(t *testing.T) {
	invocations := 0
	w := assertingWriter{t: t, msg: msgMock, onWrite: func() { invocations++ }}
	log.SetOutput(w)

	conditionallyLog(true, msgMock)
	assert.Equal(t, 1, invocations)
}

func TestConditionallyLogFalse(t *testing.T) {
	invocations := 0
	w := assertingWriter{t: t, msg: msgMock, onWrite: func() { invocations++ }}
	log.SetOutput(w)

	conditionallyLog(false, msgMock)
	assert.Equal(t, 0, invocations)
}

func mockControlPacket(syn, ack, fin, err bool) *packet.Packet {
	p, _ := packet.NewPacket(0, 0, nil)
	if syn {
		p.SetFlagSYN()
	}
	if ack {
		p.SetFlagACK()
	}
	if fin {
		p.SetFlagFIN()
	}
	if err {
		p.SetFlagERR()
	}
	return p
}
