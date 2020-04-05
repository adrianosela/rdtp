package service

import (
	"github.com/adrianosela/rdtp/netwk"
	"github.com/adrianosela/rdtp/socket"
	"github.com/pkg/errors"
)

// Service is an abstraction of the rdtp service
type Service struct {
	// network is an interface that takes packets
	// and ships them out to the network
	network *netwk.Network

	sckmgr *socket.Manager
}

// NewService returns an rdtp service instance
func NewService() (*Service, error) {
	nw, err := netwk.NewNetwork()
	if err != nil {
		return nil, errors.Wrap(err, "could not acquire network")
	}
	mgr, err := socket.NewManager()
	if err != nil {
		return nil, errors.Wrap(err, "could not initialize socket manager")
	}
	return &Service{
		network: nw,
		sckmgr:  mgr,
	}, nil
}
