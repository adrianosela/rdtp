package listener

import (
	"net"

	"github.com/adrianosela/rdtp"
	"github.com/pkg/errors"
)

// Listener represents a connection to the application listening on a given port
type Listener struct {
	Port     uint16
	notifier net.Conn
}

// New is the Listener constructor
func New(port uint16, c net.Conn) *Listener {
	return &Listener{
		Port:     port,
		notifier: c,
	}
}

// Notify sends an rdtp service message to a listener about an inbound
// connection from a remote address
func (l *Listener) Notify(connectingRemoteAddress *rdtp.Addr) error {
	msg, err := rdtp.NewServiceMessage(rdtp.ServiceMessageTypeNotify, nil, connectingRemoteAddress, nil)
	if err != nil {
		return errors.Wrap(err, "could not create NOTIFY service message")
	}
	if _, err := l.notifier.Write(msg); err != nil {
		return errors.Wrap(err, "could not write notification (packet address) to application")
	}
	return nil
}

// Close closes the listener's notifier connection
func (l *Listener) Close() error {
	return l.notifier.Close()
}
