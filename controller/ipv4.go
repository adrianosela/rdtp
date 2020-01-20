package controller

import (
	"log"
	"net"

	"github.com/adrianosela/rdtp"
	"github.com/pkg/errors"
)

const (
	ipv4HeaderProtocolOffset = byte(9)
	// unassigned as per https://en.wikipedia.org/wiki/List_of_IP_protocol_numbers
	ipv4HeaderProtocolRDTP = byte(250)
)

func (ctrl *Controller) listenRDTPOverIPv4() error {
	// listen for IP packets on all available interfaces
	conn, err := net.ListenIP("ip:ip", nil)
	if err != nil {
		return errors.Wrap(err, "could not listen for IP")
	}

	for {
		// allocate buffer for packet
		buf := make([]byte, rdtp.MaxPacketBytes+ipHeaderLen)
		// ReadFrom() on an IPConn handles stripping the IP header
		// To examine ip header values we can use Read() or ReadString()
		eof, err := conn.Read(buf)
		if err != nil {
			// soft error
			log.Println(errors.Wrap(err, "error reading IP buffer"))
			continue
		}

		if eof < ipHeaderLen+rdtp.HeaderByteSize {
			log.Println(errors.Wrap(err, "too-small IP packet"))
			continue
		}

		if []byte(buf)[ipv4HeaderProtocolOffset] != ipv4HeaderProtocolRDTP {
			// drop non rdtp packets (anything over IP comes through)
			continue
		}

		p, err := rdtp.Deserialize([]byte(buf)[ipHeaderLen:eof])
		if err != nil {
			log.Println(errors.Wrap(err, "could not build received rdtp packet"))
			continue
		}

		if err = ctrl.MultiplexPacket(p); err != nil {
			log.Println(errors.Wrap(err, "could not multiplex rdtp packet"))
			continue
		}
	}
}
