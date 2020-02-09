package service

import (
	"github.com/adrianosela/rdtp/ports"
	"github.com/adrianosela/rdtp/ports/filesystem"
	"github.com/pkg/errors"
)

// Service represents the RDTP service
type Service interface {
	// TODO (e.g. NewConn())
}

// Default is the default RDTP service
type Default struct {
	ports ports.Manager
}

// Acquire returns the default RDTP service
func Acquire() (Service, error) {
	mgr, err := filesystem.NewFSManager("")
	if err != nil {
		return nil, errors.Wrap(err, "could not init file system ports manager")
	}
	return &Default{
		ports: mgr,
	}, nil
}
