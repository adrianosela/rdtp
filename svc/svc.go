package svc

import (
	"sync"

	"github.com/adrianosela/rdtp/socket"
)

/* Architecture notes:

outbound: user =msg=> svc =msg=> connworker =pck=> atc =pck=> netwk
inbound: nwtwk =pck=> atc =pck=> connworker =msg=> svc =msg=> user

- The user is in charge of read() and write() to a net.Conn
- The svc is in charge of dispatching/killing connection worker for every
  outbound or inbound connection. The svc checks the status of the worker
  before reading from it or writing to it
- The connection worker is in charge of packetizing outgoing messages
  and de-packetizing incomming packets. Each connection worker has its own
  ATC
- The ATC is in charge of maintaing statistics about the connections
  transmission/reception, etc- as well as retrnsmissions and packet loss.
  It has a direct connection with the network. The ATC can directly talk to
  the svc in order


we will worry about listening for inbound SYNs later...
*/

// RDTPService represents the rdtp service
type RDTPService struct {
	sync.RWMutex

	// sockets is a map of sockets where each socket's
	// unique identifier is "laddr:lport raddr:rport",
	// e.g. "192.168.1.75:4444 192.168.1.88:1201"
	sockets map[string]*socket.Socket
}

// NewRDTPService returns an initialized RDTP service
func NewRDTPService() *RDTPService {
	return &RDTPService{
		sockets: make(map[string]*socket.Socket),
	}
}
