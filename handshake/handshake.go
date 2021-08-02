package handshake

import (
	"fmt"
	"log"
	"time"

	"github.com/adrianosela/rdtp/packet"
	"github.com/pkg/errors"
)

const (
	debug   = false
	flagFmt = "{SYN[%t] ACK[%t] FIN[%t] ERR[%t]}"
)

type ctrlPacketSender func(syn, ack, fin, err bool) error

// InitiateConnection sends a SYN, waits for a SYN ACK, and sends an ACK
func InitiateConnection(recv chan *packet.Packet, recvTimeout time.Duration, sendCtrl ctrlPacketSender) error {
	// send SYN
	if err := sendCtrl(true, false, false, false); err != nil {
		conditionallyLog(debug, "DIAL: Send SYN [FAIL]: %s", err)
		return errors.Wrap(err, "connect handshake failed when sending SYN")
	}
	conditionallyLog(debug, "DIAL: Send SYN [OK]")

	// wait for SYN ACK
	if err := receiveControlPacket(recv, true, true, false, false, recvTimeout); err != nil {
		conditionallyLog(debug, "DIAL: Receive SYN ACK [FAIL]: %s", err)
		return errors.Wrap(err, "connect handshake failed when waiting for SYN ACK")
	}
	conditionallyLog(debug, "DIAL: Receive SYN ACK [OK]")

	// send ACK
	if err := sendCtrl(false, true, false, false); err != nil {
		conditionallyLog(debug, "DIAL: Send ACK [FAIL]: %s", err)
		return errors.Wrap(err, "connect handshake failed when sending ACK")
	}
	conditionallyLog(debug, "DIAL: Send ACK [OK]")

	return nil
}

// AcceptConnection sends a SYN ACK and waits for an ACK
func AcceptConnection(recv chan *packet.Packet, recvTimeout time.Duration, sendCtrl ctrlPacketSender) error {
	// send SYN ACK
	if err := sendCtrl(true, true, false, false); err != nil {
		conditionallyLog(debug, "ACCEPT: Send SYN ACK [FAIL]: %s", err)
		return errors.Wrap(err, "connect handshake failed when sending SYN ACK")
	}
	conditionallyLog(debug, "ACCEPT: Send SYN ACK [OK]")

	// wait for ACK
	if err := receiveControlPacket(recv, false, true, false, false, recvTimeout); err != nil {
		conditionallyLog(debug, "ACCEPT: Receive ACK [FAIL]: %s", err)
		return errors.Wrap(err, "connect handshake failed when waiting for ACK")
	}
	conditionallyLog(debug, "ACCEPT: Receive ACK [OK]")

	return nil
}

// InitiateDisconnection sends a FIN, waits for a FIN ACK, and sends an ACK
func InitiateDisconnection(recv chan *packet.Packet, recvTimeout time.Duration, sendCtrl ctrlPacketSender) error {
	// SEND FIN
	if err := sendCtrl(false, false, true, false); err != nil {
		conditionallyLog(debug, "FINISH (closed by local): Send FIN [FAIL]: %s", err)
		return errors.Wrap(err, "finish handshake failed when sending FIN")
	}
	conditionallyLog(debug, "FINISH (closed by local): Send FIN [OK]")

	// wait for FIN ACK
	if err := receiveControlPacket(recv, false, true, true, false, recvTimeout); err != nil {
		conditionallyLog(debug, "FINISH (closed by local): Receive FIN ACK [FAIL]: %s", err)
		return errors.Wrap(err, "finish handshake failed when waiting for FIN ACK")
	}
	conditionallyLog(debug, "FINISH (closed by local): Receive FIN ACK [OK]")

	// send ACK
	if err := sendCtrl(false, true, false, false); err != nil {
		conditionallyLog(debug, "FINISH (closed by local): Send ACK [FAIL]: %s", err)
		return errors.Wrap(err, "finish handshake failed when sending ACK")
	}
	conditionallyLog(debug, "FINISH (closed by local): Send ACK [OK]")

	return nil
}

// AcceptDisconnection sends a FIN ACK and waits for an ACK
func AcceptDisconnection(recv chan *packet.Packet, recvTimeout time.Duration, sendCtrl ctrlPacketSender) error {
	// send FIN ACK
	if err := sendCtrl(false, true, true, false); err != nil {
		conditionallyLog(debug, "FINISH (closed by remote): Send FIN ACK [FAIL]: %s", err)
		return errors.Wrap(err, "finish handshake failed when sending FIN ACK")
	}
	conditionallyLog(debug, "FINISH (closed by remote): Send FIN ACK [OK]")

	// wait for ACK
	if err := receiveControlPacket(recv, false, true, false, false, recvTimeout); err != nil {
		conditionallyLog(debug, "FINISH (closed by remote): Receive ACK [FAIL]: %s", err)
		return errors.Wrap(err, "finish handshake failed when waiting for ACK")
	}
	conditionallyLog(debug, "FINISH (closed by remote): Receive ACK [OK]")

	return nil
}

// receiveControlPacket blocks until the a packet is received (or timeout)
func receiveControlPacket(in chan *packet.Packet, syn, ack, fin, err bool, recvTimeout time.Duration) error {
	for {
		select {
		case p := <-in:
			if syn != p.IsSYN() || ack != p.IsACK() || fin != p.IsFIN() || err != p.IsERR() {
				return fmt.Errorf(
					"expected packet with flags %s but got %s",
					fmt.Sprintf(flagFmt, syn, ack, fin, err),
					fmt.Sprintf(flagFmt, p.IsSYN(), p.IsACK(), p.IsFIN(), p.IsERR()))
			}
			return nil
		case <-time.After(recvTimeout):
			return errors.New("operation timed out")
		}
	}
}

func conditionallyLog(cond bool, fmtString string, indirects ...interface{}) {
	if cond {
		log.Printf(fmtString, indirects...)
	}
}
