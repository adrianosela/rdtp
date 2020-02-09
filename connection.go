package rdtp

import (
	"github.com/adrianosela/rdtp/service"
	"github.com/pkg/errors"
)

// Conn is an RDTP connection
type Conn struct {
	svcConn service.Service
}

// Dial establishes an RDTP connection with a remote IP host
func Dial(ip string) (*Conn, error) {
	svc, err := service.Acquire()
	if err != nil {
		return nil, errors.Wrap(err, "could not acquire default RDTP service")
	}
	// TODO
	return &Conn{
		svcConn: svc,
	}, nil
}

// Close closes an RDTP connection
func (c *Conn) Close() error {
	// TODO
	return nil
}
