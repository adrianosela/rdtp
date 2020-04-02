package mux

import (
	"errors"
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/adrianosela/rdtp"
	"github.com/adrianosela/rdtp/packet"
)

// MemoryMux is an in-memory implementation of a Mux
type MemoryMux struct {
	sync.RWMutex // inherit read/write lock behavior

	conns map[uint16]net.Conn // port number to client conn
}

// NewMemoryMux returns an initialized in-memory multiplexer
func NewMemoryMux() *MemoryMux {
	return &MemoryMux{
		conns: make(map[uint16]net.Conn),
	}
}

// MultiplexPacket delivers a packet to the correct destination
func (m *MemoryMux) MultiplexPacket(p *packet.Packet) error {
	info := fmt.Sprintf("[MUX] RX %d ==> %d", p.Length, p.DstPort)

	if !p.Check() {
		log.Print(fmt.Sprintf("%s %s", info, "✘ [CORRUPTED]"))
		return errors.New("checksum failed")
	}

	conn, error := m.Get(p.DstPort)
	if error != nil {
		log.Print(fmt.Sprintf("%s %s", info, "✘ [PORT CLOSED]"))
		return fmt.Errorf("port %d closed", p.DstPort)
	}

	if _, err := conn.Write(p.Payload); err != nil {
		log.Print(fmt.Sprintf("%s %s", info, "✘ [ERROR]"))
		m.Detach(p.DstPort) // close port on error
		return errors.New("could not forward payload to client")
	}

	log.Print(fmt.Sprintf("%s %s", info, "✔ [OK]"))
	return nil
}

// Get returns an indexed connection, nil if not set
func (m *MemoryMux) Get(p uint16) (net.Conn, error) {
	m.RLock()
	defer m.RUnlock()

	conn, ok := m.conns[p]
	if !ok {
		return nil, fmt.Errorf("no connection on port %d", p)
	}

	return conn, nil
}

// Attach attaches a connection to a given port
func (m *MemoryMux) Attach(p uint16, c net.Conn) error {
	m.Lock()
	defer m.Unlock()

	if _, ok := m.conns[p]; ok {
		return fmt.Errorf("port %d is in use", p)
	}

	m.conns[p] = c
	return nil
}

// AttachAny attaches a connection to the lowest unused port
func (m *MemoryMux) AttachAny(c net.Conn) (uint16, error) {
	m.Lock()
	defer m.Unlock()

	// give out first unused port
	for port := uint16(1); port < rdtp.MaxPort; port++ {
		if _, ok := m.conns[port]; !ok {
			m.conns[port] = c
			return port, nil
		}
	}

	return uint16(0), fmt.Errorf("all ports in use")
}

// Detach detaches an indexed connection from the multiplexer
func (m *MemoryMux) Detach(p uint16) {
	m.Lock()
	defer m.Unlock()
	delete(m.conns, p)
}
