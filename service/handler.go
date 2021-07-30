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
		// TODO send error message
		return
	}

	var req rdtp.ClientMessage
	if err := json.Unmarshal(buf[:n], &req); err != nil {
		// TODO send error message
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
		// TODO send error message
		break
	}

	return
}

func (s *Service) handleClientMessageDial(c net.Conn, r rdtp.ClientMessage) {
	defer c.Close()

	laddr := &rdtp.Addr{Host: getOutboundIP(), Port: uint16(rand.Intn(int(rdtp.MaxPort)-1) + 1)}
	sck, err := socket.NewSocket(socket.Config{
		LocalAddr:          laddr,
		RemoteAddr:         &r.RemoteAddr,
		ToApplicationLayer: c,
		ToController:       s.netLayer.Send,
	})
	if err != nil {
		// TODO send error message
		// ("could not get socket for user")
		return
	}

	if err = s.sckmgr.Put(sck); err != nil {
		// TODO send error message
		// ("could not attach socket to socket manager")
		return
	}
	log.Printf("%s [attached]\n", sck.ID())

	// send syn
	if err := s.sendControlPacket(laddr, &r.RemoteAddr, true, false); err != nil {
		// TODO send error message
		// ("could not send SYN control packet to remote")
		return
	}

	defer func() {
		s.sckmgr.Evict(sck.ID())
		log.Printf("%s [evicted]\n", sck.ID())
	}()

	msg, err := rdtp.NewServiceMessage(rdtp.ServiceMessageTypeOK, laddr, &r.RemoteAddr, nil)
	if err != nil {
		// TODO send error message
		// ("could not create OK message")
		return
	}

	if _, err = c.Write(msg); err != nil {
		// TODO send error message
		// ("could not reply with OK message")
		return
	}

	if err = sck.Start(); err != nil {
		// TODO send error message
		// ("socket failure")
		return
	}

	return
}

func (s *Service) handleClientMessageAccept(c net.Conn, r rdtp.ClientMessage) {
	defer c.Close()

	sck, err := socket.NewSocket(socket.Config{
		LocalAddr:          &r.LocalAddr,
		RemoteAddr:         &r.RemoteAddr,
		ToApplicationLayer: c,
		ToController:       s.netLayer.Send,
	})
	if err != nil {
		// TODO send error message
		// ("could not get socket for user")
		return
	}

	if err = s.sckmgr.Put(sck); err != nil {
		// TODO send error message
		// ("could not attach socket to socket manager")
		return
	}
	log.Printf("%s [attached]\n", sck.ID())

	defer func() {
		s.sckmgr.Evict(sck.ID())
		log.Printf("%s [evicted]\n", sck.ID())
	}()

	msg, err := rdtp.NewServiceMessage(rdtp.ServiceMessageTypeOK, &r.LocalAddr, &r.RemoteAddr, nil)
	if err != nil {
		// TODO send error message
		// ("could not create OK message")
		return
	}

	if _, err = c.Write(msg); err != nil {
		// TODO send error message
		// ("could not reply with OK message")
		return
	}

	if err = sck.Start(); err != nil {
		// TODO send error message
		// ("socket failure")
		return
	}

	return
}

func (s *Service) handleClientMessageListen(c net.Conn, r rdtp.ClientMessage) {
	if err := s.sckmgr.PutListener(socket.NewListener(r.LocalAddr.Port, c)); err != nil {
		// TODO send error message
		// (fmt.Sprintf("could not attach listener to port %d", r.LocalAddr.Port)
		return
	}
	msg, err := rdtp.NewServiceMessage(rdtp.ServiceMessageTypeOK, &r.LocalAddr, nil, nil)
	if err != nil {
		// TODO send error message
		// ("could not create OK message")
		return
	}
	if _, err = c.Write(msg); err != nil {
		// TODO send error message
		// ("could not reply with OK message")
		return
	}

	return
}

func (s *Service) sendControlPacket(laddr, raddr *rdtp.Addr, syn, ack bool) error {
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
	p.SetSourceIPv4(net.ParseIP(laddr.Host))
	p.SetDestinationIPv4(net.ParseIP(raddr.Host))
	p.SetSum()

	// send SYN to destination
	if err = s.netLayer.Send(p); err != nil {
		return errors.Wrap(err, "could not send SYN to destination")
	}

	return nil
}
