package service

import (
	"github.com/adrianosela/rdtp/ipv4"
	"github.com/adrianosela/rdtp/socket"
	"github.com/pkg/errors"
)

// Service is an abstraction of the rdtp service
type Service struct {
	network *ipv4.IPv4
	sckmgr  *socket.Manager
}

// NewService returns an rdtp service instance
func NewService() (*Service, error) {
	ip, err := ipv4.NewIPv4()
	if err != nil {
		return nil, errors.Wrap(err, "could not acquire network")
	}
	mgr, err := socket.NewManager()
	if err != nil {
		return nil, errors.Wrap(err, "could not initialize socket manager")
	}
	return &Service{
		network: ip,
		sckmgr:  mgr,
	}, nil
}

// Start starts the rdtp service
func (s *Service) Start() error {
	for {
		// TODO
	}
}
