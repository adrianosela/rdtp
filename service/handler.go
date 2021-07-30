package service

import (
	"encoding/json"
	"log"
	"math/rand"
	"net"

	"github.com/adrianosela/rdtp"
	"github.com/adrianosela/rdtp/packet"
	"github.com/adrianosela/rdtp/service/socket"
	"github.com/pkg/errors"
)

func (s *Service) handleClientMessage(c net.Conn) {
	buf := make([]byte, 1024)
	n, err := c.Read(buf)
	if err != nil {
		log.Println("connection closed by client")
		sendErrorMessage(c, rdtp.ServiceErrorTypeConnClosedByClient)
		return
	}

	var req rdtp.ClientMessage
	if err := json.Unmarshal(buf[:n], &req); err != nil {
		log.Println("malformed client message received")
		sendErrorMessage(c, rdtp.ServiceErrorTypeMalformedMessage)
		return
	}

	switch req.Type {
	case rdtp.ClientMessageTypeAccept:
		s.handleClientMessageAccept(c, req)
		break
	case rdtp.ClientMessageTypeDial:
		s.handleClientMessageDial(c, req)
		break
	case rdtp.ClientMessageTypeListen:
		s.handleClientMessageListen(c, req)
		break
	default:
		log.Println("invalid message type received")
		sendErrorMessage(c, rdtp.ServiceErrorTypeInvalidMessageType)
		break
	}

	return
}

func (s *Service) handleClientMessageDial(c net.Conn, r rdtp.ClientMessage) {
	laddr := &rdtp.Addr{Host: getOutboundIP(), Port: uint16(rand.Intn(int(rdtp.MaxPort)-1) + 1)}
	sck, err := socket.NewSocket(socket.Config{
		LocalAddr:          laddr,
		RemoteAddr:         &r.RemoteAddr,
		ToApplicationLayer: c,
		ToController:       s.netLayer.Send,
	})
	if err != nil {
		c.Close()
		log.Println(errors.Wrap(err, "failed to create socket"))
		sendErrorMessage(c, rdtp.ServiceErrorTypeFailedToCreateSocket)
		return
	}

	if err = s.sckmgr.Put(sck); err != nil {
		sck.Close()
		log.Println(errors.Wrap(err, "failed to attach socket"))
		sendErrorMessage(c, rdtp.ServiceErrorTypeFailedToAttachSocket)
		return
	}
	defer s.sckmgr.Evict(sck.ID())

	// send syn
	if err := s.sendControlPacket(laddr, &r.RemoteAddr, true, false, false); err != nil {
		log.Println(errors.Wrap(err, "handshake failed"))
		sendErrorMessage(c, rdtp.ServiceErrorTypeFailedHandshake)
		return
	}

	if err := sendOKMessage(c, laddr, &r.RemoteAddr); err != nil {
		log.Println(errors.Wrap(err, "failed to send ok message"))
		sendErrorMessage(c, rdtp.ServiceErrorTypeFailedCommunication)
		return
	}

	if err = sck.Run(); err != nil {
		log.Println(errors.Wrap(err, "socket run failed"))
		sendErrorMessage(c, rdtp.ServiceErrorTypeFailedSocketRun)
		return
	}

	// send fin, dont care about error
	s.sendControlPacket(laddr, &r.RemoteAddr, false, false, true)
	return
}

func (s *Service) handleClientMessageAccept(c net.Conn, r rdtp.ClientMessage) {
	sck, err := socket.NewSocket(socket.Config{
		LocalAddr:          &r.LocalAddr,
		RemoteAddr:         &r.RemoteAddr,
		ToApplicationLayer: c,
		ToController:       s.netLayer.Send,
	})
	if err != nil {
		c.Close()
		log.Println(errors.Wrap(err, "failed to create socket"))
		sendErrorMessage(c, rdtp.ServiceErrorTypeFailedToCreateSocket)
		return
	}

	if err = s.sckmgr.Put(sck); err != nil {
		sck.Close()
		log.Println(errors.Wrap(err, "failed to attach socket"))
		sendErrorMessage(c, rdtp.ServiceErrorTypeFailedToAttachSocket)
		return
	}
	defer s.sckmgr.Evict(sck.ID())

	if err := sendOKMessage(c, &r.LocalAddr, &r.RemoteAddr); err != nil {
		log.Println(errors.Wrap(err, "failed to send ok message"))
		sendErrorMessage(c, rdtp.ServiceErrorTypeFailedCommunication)
		return
	}

	if err = sck.Run(); err != nil {
		log.Println(errors.Wrap(err, "socket run failed"))
		sendErrorMessage(c, rdtp.ServiceErrorTypeFailedSocketRun)
		return
	}

	// send fin, dont care about error
	s.sendControlPacket(&r.LocalAddr, &r.RemoteAddr, false, false, true)
	return
}

func (s *Service) handleClientMessageListen(c net.Conn, r rdtp.ClientMessage) {
	if err := s.sckmgr.PutListener(socket.NewListener(r.LocalAddr.Port, c)); err != nil {
		log.Println(errors.Wrap(err, "failed to attach listener"))
		sendErrorMessage(c, rdtp.ServiceErrorTypeFailedToAttachListener)
		return
	}
	defer s.sckmgr.EvictListener(r.LocalAddr.Port)

	if err := sendOKMessage(c, &r.LocalAddr, &r.RemoteAddr); err != nil {
		log.Println(errors.Wrap(err, "failed to send ok message"))
		sendErrorMessage(c, rdtp.ServiceErrorTypeFailedCommunication)
		return
	}

	// wait for EOF (connection close by client)
	c.Read(make([]byte, 1))

	return
}

func (s *Service) sendControlPacket(laddr, raddr *rdtp.Addr, syn, ack, fin bool) error {
	p, err := packet.NewPacket(laddr.Port, raddr.Port, nil)
	if err != nil {
		return errors.Wrap(err, "could not create new packet")
	}
	if syn {
		p.SetFlagSYN()
	}
	if ack {
		p.SetFlagACK()
	}
	if fin {
		p.SetFlagFIN()
	}
	p.SetSourceIPv4(net.ParseIP(laddr.Host))
	p.SetDestinationIPv4(net.ParseIP(raddr.Host))
	p.SetSum()

	// send SYN to destination
	if err = s.netLayer.Send(p); err != nil {
		return errors.Wrap(err, "could not send SYN to destination")
	}

	return nil
}
