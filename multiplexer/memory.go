package multiplexer

import (
	"errors"
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/adrianosela/rdtp/packet"
)

// MapMux is an in-memory implementation of a mux
type MapMux struct {
	sync.RWMutex // inherit read/write lock behavior

	conns map[uint16]net.Conn // port number to client conn
}

// NewMapMux returns an initialized in-memory multiplexer
func NewMapMux() *MapMux {
	return &MapMux{
		conns: make(map[uint16]net.Conn),
	}
}

// MultiplexPacket delivers a packet to the correct destination
func (m *MapMux) MultiplexPacket(p *packet.Packet) error {
	info := fmt.Sprintf("[MUX] RX %d ==> %d", p.Length, p.DstPort)

	if !p.Check() {
		log.Print(fmt.Sprintf("%s %s", info, "✘ [CORRUPTED]"))
		return errors.New("checksum failed")
	}

	conn, ok := m.Get(p.DstPort)
	if !ok {
		log.Print(fmt.Sprintf("%s %s", info, "✘ [PORT CLOSED]"))
		return fmt.Errorf("port %d closed", p.DstPort)
	}

	if _, err := conn.Write(p.Payload); err != nil {
		log.Print(fmt.Sprintf("%s %s", info, "✘ [ERROR]"))
		return errors.New("could not forward payload to client")
	}

	log.Print(fmt.Sprintf("%s %s", info, "✔ [OK]"))
	return nil
}

// Get returns an indexed connection, nil if not set
func (m *MapMux) Get(p uint16) (net.Conn, bool) {
	m.RLock()
	defer m.RUnlock()
	conn, ok := m.conns[p]
	return conn, ok
}

// Attach attaches an indexed connection to the multiplexer
func (m *MapMux) Attach(p uint16, c net.Conn) {
	m.Lock()
	defer m.Unlock()
	m.conns[p] = c
}

// Detach detaches an indexed connection from the multiplexer
func (m *MapMux) Detach(p uint16) {
	m.Lock()
	defer m.Unlock()
	delete(m.conns, p)
}
