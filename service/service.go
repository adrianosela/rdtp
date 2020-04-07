package service

import (
	"github.com/adrianosela/rdtp/atc"
	"github.com/adrianosela/rdtp/ipv4"
	"github.com/adrianosela/rdtp/packet"
	"github.com/adrianosela/rdtp/socket"
	"github.com/pkg/errors"
)

// Service is an abstraction of the rdtp service
type Service struct {
	network *ipv4.IPv4
	atc     *atc.AirTrafficCtrl
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

	atc.NewAirTrafficCtrl(func(p *packet.Packet) error {
		/* TODO */
		return nil
	})

	return &Service{
		network: ip,
		sckmgr:  mgr,
	}, nil
}

// Start starts the rdtp service
func (s *Service) Start() error {
	// forward all rdtp packets received on the network to their socket
	go s.network.ForwardRDTP(func(p *packet.Packet) error {
		if err := s.sckmgr.Deliver(p); err != nil {
			return errors.Wrap(err, "could not deliver packet to rdtp socket")
		}
		return nil
	})

	for {
		/* TODO */
	}
}
