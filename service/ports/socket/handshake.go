package socket

import (
	"fmt"
	"log"
	"time"

	"github.com/adrianosela/rdtp/packet"
	"github.com/pkg/errors"
)

const (
	flagFmt                  = "{SYN[%t] ACK[%t] FIN[%t] ERR[%t]}"
	controlPacketWaitTimeout = time.Second * 1
)

// Dial sends a SYN, waits for a SYN ACK, and sends an ACK
func (s *Socket) Dial() error {
	// send SYN
	if err := s.packetizer.SendControlPacket(true, false, false, false); err != nil {
		// log.Printf("DIAL: Send SYN [FAIL]: %s", err)
		return errors.Wrap(err, "connect handshake failed when sending SYN")
	}
	// log.Println("DIAL: Send SYN [OK]")

	// wait for SYN ACK
	if err := receiveControlPacket(s.inbound, true, true, false, false, controlPacketWaitTimeout); err != nil {
		// log.Printf("DIAL: Receive SYN ACK [FAIL]: %s", err)
		return errors.Wrap(err, "connect handshake failed when waiting for SYN ACK")
	}
	// log.Println("DIAL: Receive SYN ACK [OK]")

	// send ACK
	if err := s.packetizer.SendControlPacket(false, true, false, false); err != nil {
		// log.Printf("DIAL: Send ACK [FAIL]: %s", err)
		return errors.Wrap(err, "connect handshake failed when sending ACK")
	}
	// log.Println("DIAL: Send ACK [OK]")

	return nil
}

// Accept sends a SYN ACK and waits for an ACK
func (s *Socket) Accept() error {
	// send SYN ACK
	if err := s.packetizer.SendControlPacket(true, true, false, false); err != nil {
		// log.Printf("ACCEPT: Send SYN ACK [FAIL]: %s", err)
		return errors.Wrap(err, "connect handshake failed when sending SYN ACK")
	}
	// log.Println("ACCEPT: Send SYN ACK [OK]")

	// wait for ACK
	if err := receiveControlPacket(s.inbound, false, true, false, false, controlPacketWaitTimeout); err != nil {
		// log.Printf("ACCEPT: Receive ACK [FAIL]: %s", err)
		return errors.Wrap(err, "connect handshake failed when waiting for ACK")
	}
	// log.Println("ACCEPT: Receive ACK [OK]")

	return nil
}

// Finish manages the finish handshake
func (s *Socket) Finish(closedByRemote bool) error {
	if closedByRemote {
		// send FIN ACK
		if err := s.packetizer.SendControlPacket(false, true, true, false); err != nil {
			log.Printf("FINISH (closed by remote): Send FIN ACK [FAIL]: %s", err)
			return errors.Wrap(err, "finish handshake failed when sending FIN ACK")
		}
		log.Println("FINISH (closed by remote): Send FIN ACK [OK]")

		// wait for ACK
		if err := receiveControlPacket(s.inbound, false, true, false, false, controlPacketWaitTimeout); err != nil {
			log.Printf("FINISH (closed by remote): Receive ACK [FAIL]: %s", err)
			return errors.Wrap(err, "finish handshake failed when waiting for ACK")
		}
		log.Println("FINISH (closed by remote): Receive ACK [OK]")

		return nil
	}
	// send FIN
	if err := s.packetizer.SendControlPacket(false, false, true, false); err != nil {
		log.Printf("FINISH (closed by local): Send FIN [FAIL]: %s", err)
		return errors.Wrap(err, "finish handshake failed when sending FIN")
	}
	log.Println("FINISH (closed by local): Send FIN [OK]")

	// wait for FIN ACK
	if err := receiveControlPacket(s.inbound, false, true, true, false, controlPacketWaitTimeout); err != nil {
		log.Printf("FINISH (closed by local): Receive FIN ACK [FAIL]: %s", err)
		return errors.Wrap(err, "finish handshake failed when waiting for FIN ACK")
	}
	log.Println("FINISH (closed by local): Receive FIN ACK [OK]")

	// send ACK
	if err := s.packetizer.SendControlPacket(false, true, false, false); err != nil {
		log.Printf("FINISH (closed by local): Send ACK [FAIL]: %s", err)
		return errors.Wrap(err, "finish handshake failed when sending ACK")
	}
	log.Println("FINISH (closed by local): Send ACK [OK]")

	return nil
}

// receiveControlPacket blocks until the a packet is received (or timeout)
func receiveControlPacket(in chan *packet.Packet, syn, ack, fin, err bool, timeout time.Duration) error {
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
		case <-time.After(timeout):
			return errors.New("operation timed out")
		}
	}
}
