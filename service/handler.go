package service

import (
	"encoding/json"
	"io"
	"log"
	"math/rand"
	"net"

	"github.com/adrianosela/rdtp"
	"github.com/adrianosela/rdtp/service/ports/listener"
	"github.com/adrianosela/rdtp/socket"
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
	sck, err := socket.New(socket.Config{
		LocalAddr:   laddr,
		RemoteAddr:  &r.RemoteAddr,
		Application: c,
		Network:     s.network,
	})
	if err != nil {
		c.Close()
		log.Println(errors.Wrap(err, "failed to create socket"))
		sendErrorMessage(c, rdtp.ServiceErrorTypeFailedToCreateSocket)
		return
	}

	if err = s.ports.Put(sck); err != nil {
		sck.Close()
		log.Println(errors.Wrap(err, "failed to attach socket"))
		sendErrorMessage(c, rdtp.ServiceErrorTypeFailedToAttachSocket)
		return
	}
	defer s.ports.Evict(sck.ID())

	if err := sck.Dial(); err != nil {
		log.Println(errors.Wrap(err, "socket dial failed"))
		sendErrorMessage(c, rdtp.ServiceErrorTypeFailedHandshake)
		return
	}

	if err := sendOKMessage(c, laddr, &r.RemoteAddr); err != nil {
		log.Println(errors.Wrap(err, "failed to send ok message"))
		sendErrorMessage(c, rdtp.ServiceErrorTypeFailedCommunication)
		return
	}

	sck.Run()

	return
}

func (s *Service) handleClientMessageAccept(c net.Conn, r rdtp.ClientMessage) {
	sck, err := socket.New(socket.Config{
		LocalAddr:   &r.LocalAddr,
		RemoteAddr:  &r.RemoteAddr,
		Application: c,
		Network:     s.network,
	})
	if err != nil {
		c.Close()
		log.Println(errors.Wrap(err, "failed to create socket"))
		sendErrorMessage(c, rdtp.ServiceErrorTypeFailedToCreateSocket)
		return
	}

	if err = s.ports.Put(sck); err != nil {
		sck.Close()
		log.Println(errors.Wrap(err, "failed to attach socket"))
		sendErrorMessage(c, rdtp.ServiceErrorTypeFailedToAttachSocket)
		return
	}
	defer s.ports.Evict(sck.ID())

	if err := sck.Accept(); err != nil {
		log.Println(errors.Wrap(err, "socket accept failed"))
		sendErrorMessage(c, rdtp.ServiceErrorTypeFailedHandshake)
		return
	}

	if err := sendOKMessage(c, &r.LocalAddr, &r.RemoteAddr); err != nil {
		log.Println(errors.Wrap(err, "failed to send ok message"))
		sendErrorMessage(c, rdtp.ServiceErrorTypeFailedCommunication)
		return
	}

	sck.Run()

	return
}

func (s *Service) handleClientMessageListen(c net.Conn, r rdtp.ClientMessage) {
	if err := s.ports.AttachListener(listener.New(r.LocalAddr.Port, c)); err != nil {
		log.Println(errors.Wrap(err, "failed to attach listener"))
		sendErrorMessage(c, rdtp.ServiceErrorTypeFailedToAttachListener)
		return
	}
	defer s.ports.DetachListener(r.LocalAddr.Port)

	if err := sendOKMessage(c, &r.LocalAddr, &r.RemoteAddr); err != nil {
		log.Println(errors.Wrap(err, "failed to send ok message"))
		sendErrorMessage(c, rdtp.ServiceErrorTypeFailedCommunication)
		return
	}

	// wait for EOF (connection close by client)
	for {
		if _, err := c.Read(make([]byte, 1)); err != nil {
			if err == io.EOF {
				return
			}
		}
	}
}
