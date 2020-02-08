package controller

import (
	"net"
	"syscall"

	"github.com/adrianosela/rdtp/proto"
	"github.com/pkg/errors"
)

// Worker represents a client for the RDTP controller
type Worker struct {
	Port  uint16 // local port
	rPort uint16 // remote port

	rAddr *syscall.SockaddrInet4 // remote address

	socket int // socket file descriptor

	rxChan chan []byte
}

// NewWorker returns an RDTP Worker struct
func NewWorker(ip string) (*Worker, error) {
	ipByte := net.ParseIP(ip)
	if ipByte == nil || len(ipByte) > 4 {
		return nil, errors.New("invalid IPv4 address")
	}

	addr := &syscall.SockaddrInet4{}
	copy(addr.Addr[:], ipByte)

	// get raw network socket (AF_INET = IPv4) to send messages on
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, proto.IPProtoRDTP)
	if err != nil {
		return nil, errors.Wrap(err, "worker could not get raw network socket")
	}

	// TODO: register worker with controller to get port

	w := &Worker{
		Port:  15, // get from rdtp controller
		rPort: proto.DiscoveryPort,
		rAddr: addr,

		socket: fd,
		rxChan: make(chan []byte),
	}

	w.syn()

	return w, nil
}

func (w *Worker) Read() ([]byte, error) {
	message, ok := <-w.rxChan
	if !ok {
		return nil, errors.New("rdtp controller closed client connection")
	}
	return message, nil
}

// Kill shuts down a worker
func (w *Worker) Kill() error {

	// TODO: deregister worker with controller to release port

	close(w.rxChan)
	return nil
}
