package ports

import (
	"fmt"
	"log"
	"sync"

	"github.com/adrianosela/rdtp"
	"github.com/adrianosela/rdtp/packet"
	"github.com/adrianosela/rdtp/service/ports/listener"

	"github.com/adrianosela/rdtp/service/ports/socket"
	"github.com/pkg/errors"
)

// Manager represents the rdtp ports manager.
// It allocates and deallocates rdtp sockets and listeners.
type Manager struct {
	sync.RWMutex

	// listeners is a map of port number to listener
	listeners map[uint16]*listener.Listener

	// sockets is a map of sockets where each socket's
	// unique identifier is "laddr:lport raddr:rport",
	// e.g. "192.168.1.75:4444 192.168.1.88:1201"
	sockets map[string]*socket.Socket
}

// NewManager returns an initialized rdtp sockets manager
func NewManager() (*Manager, error) {
	return &Manager{
		listeners: make(map[uint16]*listener.Listener),
		sockets:   make(map[string]*socket.Socket),
	}, nil
}

// Get gets a socket given its id
func (m *Manager) Get(id string) (*socket.Socket, error) {
	m.RLock()
	defer m.RUnlock()

	s, ok := m.sockets[id]
	if !ok {
		return nil, errors.New("socket address not active")
	}
	return s, nil
}

// Put attaches a socket to the manager
func (m *Manager) Put(s *socket.Socket) error {
	m.Lock()
	defer m.Unlock()

	id := s.ID()
	if _, ok := m.sockets[id]; ok {
		return errors.New("socket address already in use")
	}
	m.sockets[id] = s

	log.Printf("%s [attached]\n", id)
	return nil
}

// Evict removes a socket given its id
func (m *Manager) Evict(id string) error {
	m.Lock()
	defer m.Unlock()

	sck, ok := m.sockets[id]
	if !ok {
		return nil // already not present
	}

	delete(m.sockets, id)

	sck.Close()

	log.Printf("%s [evicted]\n", id)
	return nil
}

// AttachListener attaches a listener to a port
func (m *Manager) AttachListener(l *listener.Listener) error {
	m.RLock()
	_, ok := m.listeners[l.Port]
	m.RUnlock()
	if ok {
		return fmt.Errorf("port %d is in use", l.Port)
	}

	m.Lock()
	defer m.Unlock()
	m.listeners[l.Port] = l

	log.Printf("listener on :%d [started]\n", l.Port)

	return nil
}

// DetachListener detaches a listener from a port
func (m *Manager) DetachListener(port uint16) error {
	m.Lock()
	defer m.Unlock()

	if l, ok := m.listeners[port]; ok {
		l.Close()
		delete(m.listeners, port)
	}

	log.Printf("listener on :%d [shutdown]\n", port)

	return nil
}

// notifyListener notifies a listener of an inbound remote connection
func (m *Manager) notifyListener(p *packet.Packet) error {
	m.RLock()
	l, ok := m.listeners[p.DstPort]
	m.RUnlock()
	if !ok {
		return fmt.Errorf("no listener on port %d", p.DstPort)
	}

	remoteAddress, err := p.GetSourceIPv4()
	if err != nil {
		return errors.Wrap(err, "could not get destination address from packet")
	}

	if err = l.Notify(&rdtp.Addr{Host: remoteAddress.String(), Port: p.SrcPort}); err != nil {
		return errors.Wrap(err, fmt.Sprintf("could not notify listener of connection from %s", remoteAddress.String()))
	}

	return nil
}

// Deliver delivers an inbound rdtp packet
func (m *Manager) Deliver(p *packet.Packet) error {
	if p.IsSYN() {
		if err := m.notifyListener(p); err != nil {
			return errors.Wrap(err, "could not notify listener")
		}
		return nil
	}

	id, err := socketIDFromPacket(p)
	if err != nil {
		return errors.Wrap(err, "could not build socket address from packet data")
	}
	if p.IsFIN() {
		return m.Evict(id)
	}

	m.RLock()
	s, ok := m.sockets[id]
	m.RUnlock()
	if !ok {
		return errors.New("socket address not active")
	}

	s.Deliver(p)
	return nil
}

func socketIDFromPacket(p *packet.Packet) (string, error) {
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
