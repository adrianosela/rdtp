package service

import (
	"log"
	"net"

	"github.com/adrianosela/rdtp"
	"github.com/pkg/errors"
)

// get preferred outbound ip of this machine
func getOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}

func sendOKMessage(c net.Conn, laddr, raddr *rdtp.Addr) error {
	msg, err := rdtp.NewServiceMessage(rdtp.ServiceMessageTypeOK, laddr, raddr, nil)
	if err != nil {
		return errors.Wrap(err, "failed to create new OK service message")
	}
	if _, err = c.Write(msg); err != nil {
		return errors.Wrap(err, "failed to send new OK service message")
	}
	return nil
}

func sendErrorMessage(c net.Conn, errType rdtp.ServiceErrorType) {
	msg, err := rdtp.NewServiceMessage(rdtp.ServiceMessageTypeError, nil, nil, &errType)
	if err != nil {
		log.Println(errors.Wrap(err, "failed to create new ERROR service message"))
		return // TODO: unrecoverable?
	}
	if _, err = c.Write(msg); err != nil {
		log.Println(errors.Wrap(err, "failed to send new ERROR service message"))
		return // TODO: unrecoverable?
	}
}
