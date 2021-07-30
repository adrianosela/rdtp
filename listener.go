package rdtp

import (
	"encoding/json"
	"fmt"
	"io"
	"net"

	"github.com/pkg/errors"
)

// Listener listens for new inbound rdtp
// connections on a local rdtp port
// Implements the net.Listener interface
// https://golang.org/pkg/net/#Listener
type Listener struct {
	laddr *Addr
	svc   net.Conn
}

// Listen announces on the local network address
func Listen(address string) (net.Listener, error) {
	svc, err := net.Dial("unix", DefaultRDTPServiceAddr)
	if err != nil {
		return nil, errors.Wrap(err, "could not connect to rdtp service")
	}

	laddr, err := fromString(address)
	if err != nil {
		return nil, errors.Wrap(err, "address is not a valid rdtp address")
	}

	req, err := NewClientMessage(ClientMessageTypeListen, laddr, nil)
	if err != nil {
		return nil, errors.Wrap(err, "could not create listen request for rdtp service")
	}

	if _, err = svc.Write(req); err != nil {
		return nil, errors.Wrap(err, "could not send listen request to rdtp service")
	}

	verifiedLocalAddr, err := waitForServiceMessageOK(svc)
	if err != nil {
		if err == io.EOF {
			return nil, errors.New("Listener terminated by rdtp service")
		}
		return nil, errors.Wrap(err, "could not receive OK message from service")
	}

	l := &Listener{
		laddr: verifiedLocalAddr,
		svc:   svc,
	}

	return l, nil
}

// Accept waits for and returns the next connection to the listener.
func (l *Listener) Accept() (net.Conn, error) {
	verifiedRemoteAddr, err := waitForServiceMessageNotify(l.svc)
	if err != nil {
		if err == io.EOF {
			return nil, errors.New("Listener terminated by rdtp service")
		}
		return nil, errors.Wrap(err, "remote address is not valid")
	}

	svc, err := net.Dial("unix", DefaultRDTPServiceAddr)
	if err != nil {
		return nil, errors.Wrap(err, "could not connect to rdtp service")
	}

	req, err := NewClientMessage(ClientMessageTypeAccept, l.laddr, verifiedRemoteAddr)
	if err != nil {
		return nil, errors.Wrap(err, "could not create accept request for rdtp service")
	}

	if _, err = svc.Write(req); err != nil {
		return nil, errors.Wrap(err, "could not send accept request to rdtp service")
	}

	verifiedLocalAddr, err := waitForServiceMessageOK(svc)
	if err != nil {
		if err == io.EOF {
			return nil, errors.New("Listener terminated by rdtp service")
		}
		return nil, errors.Wrap(err, "could not receive OK message from service")
	}

	return &Conn{
		laddr: verifiedLocalAddr,
		raddr: verifiedRemoteAddr,
		svc:   svc,
	}, nil
}

// Close closes the listener.
func (l *Listener) Close() error {
	return l.svc.Close()
}

// Addr returns the listener's network address.
func (l *Listener) Addr() net.Addr {
	return l.laddr
}

func waitForServiceMessageOK(c net.Conn) (*Addr, error) {
	buf := make([]byte, 1024)
	n, err := c.Read(buf)
	if err != nil {
		if err == io.EOF {
			return nil, err
		}
		return nil, errors.Wrap(err, "error reading rdtp service message")
	}

	var msg ServiceMessage
	if err := json.Unmarshal(buf[:n], &msg); err != nil {
		return nil, errors.Wrap(err, "invalid request json")
	}

	if msg.Type == ServiceMessageTypeError {
		return nil, fmt.Errorf("Error service message: %s", msg.Error)
	}

	if msg.Type != ServiceMessageTypeOK {
		return nil, fmt.Errorf("Not OK service message type %s", msg.Type)
	}

	return &msg.LocalAddr, nil
}

func waitForServiceMessageNotify(c net.Conn) (*Addr, error) {
	buf := make([]byte, 1024)
	n, err := c.Read(buf)
	if err != nil {
		if err == io.EOF {
			return nil, err
		}
		return nil, errors.Wrap(err, "error reading rdtp service message")
	}

	var msg ServiceMessage
	if err := json.Unmarshal(buf[:n], &msg); err != nil {
		return nil, errors.Wrap(err, "invalid request json")
	}

	if msg.Type == ServiceMessageTypeError {
		return nil, fmt.Errorf("Error service message: %s", msg.Error)
	}

	if msg.Type != ServiceMessageTypeNotify {
		return nil, fmt.Errorf("Not NOTIFY service message type %s", msg.Type)
	}

	return &msg.RemoteAddr, nil
}
