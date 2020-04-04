package ctrl

import (
	"sync"

	"github.com/adrianosela/rdtp/netwk"
	"github.com/adrianosela/rdtp/socket"
	"github.com/pkg/errors"
)

// Controller represents the rdtp controller.
//
// It essentially foresees the following movement of data:
// - outbound: user =msg=> socket =pck=> atc =pck=> netwk
// - inbound: nwtwk =pck=> atc =pck=> socket =msg=> user
type Controller struct {
	sync.RWMutex

	// network is an interface that takes packets
	// and ships them out to the network
	network *netwk.Network

	// sockets is a map of sockets where each socket's
	// unique identifier is "laddr:lport raddr:rport",
	// e.g. "192.168.1.75:4444 192.168.1.88:1201"
	sockets map[string]*socket.Socket
}

// NewController returns an initialized rdtp controller
func NewController() (*Controller, error) {
	nw, err := netwk.NewNetwork()
	if err != nil {
		return nil, errors.Wrap(err, "could not attach controller to network")
	}
	return &Controller{
		network: nw,
		sockets: make(map[string]*socket.Socket),
	}, nil
}

// Get gets a socket given its id
func (ctrl *Controller) Get(id string) (*socket.Socket, error) {
	ctrl.RLock()
	defer ctrl.RUnlock()

	s, ok := ctrl.sockets[id]
	if !ok {
		return nil, errors.New("socket address not active")
	}
	return s, nil
}

// Put attaches a socket to the controller
func (ctrl *Controller) Put(s *socket.Socket) error {
	ctrl.Lock()
	defer ctrl.Unlock()

	id := s.ID()
	if _, ok := ctrl.sockets[id]; ok {
		return errors.New("socket address already in use")
	}
	ctrl.sockets[id] = s
	return nil
}

// Evict removes a socket given its id
func (ctrl *Controller) Evict(id string) error {
	ctrl.Lock()
	defer ctrl.Unlock()

	if _, ok := ctrl.sockets[id]; !ok {
		return errors.New("socket address not active")
	}
	delete(ctrl.sockets, id)
	return nil
}
