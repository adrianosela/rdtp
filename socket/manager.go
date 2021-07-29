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

	// listeners is a map of port number to listener
	listeners map[uint16]*Listener

	// sockets is a map of sockets where each socket's
	// unique identifier is "laddr:lport raddr:rport",
	// e.g. "192.168.1.75:4444 192.168.1.88:1201"
	sockets map[string]*Socket
}

// NewManager returns an initialized rdtp sockets manager
func NewManager() (*Manager, error) {
	return &Manager{
		listeners: make(map[uint16]*Listener),
		sockets:   make(map[string]*Socket),
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
		return errors.New("socket address not active")
	}

	s.inbound <- p
	return nil
}

// PutListener attaches a listener to a port
func (m *Manager) PutListener(port uint16, l *Listener) error {
	m.RLock()
	_, ok := m.listeners[port]
	m.RUnlock()
	if ok {
		return fmt.Errorf("port %d is in use", port)
	}

	m.Lock()
	defer m.Unlock()
	m.listeners[port] = l

	return nil
}

// EvictListener detaches a listener from a port
func (m *Manager) EvictListener(port uint16) error {
	m.Lock()
	defer m.Unlock()

	if l, ok := m.listeners[port]; ok {
		l.application.Close()
		delete(m.listeners, port)
	}

	return nil
}

// NotifyListener notifies a listener of a SYN packet
func (m *Manager) NotifyListener(p *packet.Packet) error {
	m.RLock()
	l, ok := m.listeners[p.DstPort]
	m.RUnlock()
	if !ok {
		return fmt.Errorf("no listener on port %d", p.DstPort)
	}

	src, err := p.GetSourceIPv4()
	if err != nil {
		return errors.Wrap(err, "could not get destination address:port from packet")
	}

	if _, err := l.application.Write([]byte(fmt.Sprintf("%s:%d", src, p.SrcPort))); err != nil {
		return errors.Wrap(err, "could not write packet address to application")
	}

	return nil
}

func idFromPacket(p *packet.Packet) (string, error) {
	// destination = local for inbound, remote for outbound pcks
	dst, err := p.GetDestinationIPv4()
	if err != nil {
		return "", err
	}
	// source = local for outbound, remote for inbound pcks
	src, err := p.GetSourceIPv4()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s:%d %s:%d", dst, p.DstPort, src, p.SrcPort), nil
}
