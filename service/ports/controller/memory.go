package controller

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

// MemoryController represents an in-memory rdtp ports manager.
// It allocates and deallocates rdtp sockets and listeners.
type MemoryController struct {
	sync.RWMutex

	// listeners is a map of port number to listener
	listeners map[uint16]*listener.Listener

	// sockets is a map of sockets where each socket's
	// unique identifier is "laddr:lport raddr:rport",
	// e.g. "192.168.1.75:4444 192.168.1.88:1201"
	sockets map[string]*socket.Socket
}

// NewMemoryController returns an initialized in-memory rdtp sockets manager
func NewMemoryController() *MemoryController {
	return &MemoryController{
		listeners: make(map[uint16]*listener.Listener),
		sockets:   make(map[string]*socket.Socket),
	}
}

// Put attaches a socket to the controller
func (m *MemoryController) Put(s *socket.Socket) error {
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
func (m *MemoryController) Evict(id string) error {
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
func (m *MemoryController) AttachListener(l *listener.Listener) error {
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
func (m *MemoryController) DetachListener(port uint16) error {
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
func (m *MemoryController) notifyListener(p *packet.Packet) error {
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
func (m *MemoryController) Deliver(p *packet.Packet) error {
	if p.IsSYN() && !p.IsACK() {
		if err := m.notifyListener(p); err != nil {
			return errors.Wrap(err, "could not notify listener")
		}
		return nil
	}

	id, err := socketIDFromPacket(p)
	if err != nil {
		return errors.Wrap(err, "could not build socket address from packet data")
	}
	if p.IsFIN() && !p.IsACK() {
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
