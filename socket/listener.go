package socket

import "net"

type Listener struct {
	application net.Conn
}

func NewListener(c net.Conn) *Listener {
	return &Listener{
		application: c,
	}
}
