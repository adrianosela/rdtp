package socket

import (
	"fmt"
	"sync"

	"github.com/adrianosela/rdtp/packet"
	"github.com/pkg/errors"
)

// Manager represents the rdtp sockets controller.
// It allocates and deallocates rdtp sockets.
type Manager struct {
	sync.RWMutex

	// sockets is a map of sockets where each socket's
	// unique identifier is "laddr:lport raddr:rport",
	// e.g. "192.168.1.75:4444 192.168.1.88:1201"
	sockets map[string]*Socket
}

// NewManager returns an initialized rdtp sockets manager
func NewManager() (*Manager, error) {
	return &Manager{
		sockets: make(map[string]*Socket),
	}, nil
}

// Get gets a socket given its id
func (m *Manager) Get(id string) (*Socket, error) {
	m.RLock()
	defer m.RUnlock()

	s, ok := m.sockets[id]
	if !ok {
		return nil, errors.New("socket address not active")
	}
	return s, nil
}

// Put attaches a socket to the manager
func (m *Manager) Put(s *Socket) error {
	m.Lock()
	defer m.Unlock()

	id := s.ID()
	if _, ok := m.sockets[id]; ok {
		return errors.New("socket address already in use")
	}
	m.sockets[id] = s
	return nil
}

// Evict removes a socket given its id
func (m *Manager) Evict(id string) error {
	m.Lock()
	defer m.Unlock()

	if _, ok := m.sockets[id]; !ok {
		return errors.New("socket address not active")
	}
	delete(m.sockets, id)
	return nil
}

// Deliver delivers an inbound rdtp packet
func (m *Manager) Deliver(p *packet.Packet) error {
	id, err := idFromPacket(p)
	if err != nil {
		return errors.Wrap(err, "could not build socket address from packet data")
	}

	m.RLock()
	s, ok := m.sockets[id]
	m.RUnlock()

	if !ok {
		errors.New("socket address not active")
	}

	s.inbound <- p
	return nil
}

func idFromPacket(p *packet.Packet) (string, error) {
	ipv4, err := p.IPv4Details()
	if err != nil {
		return "", errors.Wrap(err, "no ipv4 data in rdtp packet")
	}
	return fmt.Sprintf("%s:%d %s:%d",
		ipv4.DstIP.String(), p.DstPort,
		ipv4.SrcIP.String(), p.SrcPort), nil
}
