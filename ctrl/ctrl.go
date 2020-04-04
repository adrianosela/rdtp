package ctrl

import (
	"sync"

	"github.com/adrianosela/rdtp/socket"
)

// Controller represents the rdtp controller. It essentially
// foresees the following movement of data:
// - outbound: user =msg=> socket =pck=> atc =pck=> netwk
// - inbound: nwtwk =pck=> atc =pck=> socket =msg=> user
type Controller struct {
	sync.RWMutex

	// sockets is a map of sockets where each socket's
	// unique identifier is "laddr:lport raddr:rport",
	// e.g. "192.168.1.75:4444 192.168.1.88:1201"
	sockets map[string]*socket.Socket
}

// NewController returns an initialized rdtp controller
func NewController() *Controller {
	return &Controller{
		sockets: make(map[string]*socket.Socket),
	}
}
