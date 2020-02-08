package controller

import (
	"syscall"

	"github.com/adrianosela/rdtp/packet"
	"github.com/pkg/errors"
)

func (w *Worker) syn() error {
	p, err := packet.NewPacket(w.Port, w.rPort, nil)
	if err != nil {
		return errors.Wrap(err, "could not build rdtp SYN packet for sending")
	}
	p.SetSYN()
	if err = syscall.Sendto(w.socket, p.Serialize(), 0, w.rAddr); err != nil {
		errors.Wrap(err, "could not send syn to network socket")
	}
	return nil
}
