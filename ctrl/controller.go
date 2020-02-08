package ctrl

import "github.com/adrianosela/rdtp/ctrl/filesystem"

// Controller is the interface in charge of managing
// ports for a given rdtp implementation
type Controller interface {
	AllocateAny() (uint16, error)
	Allocate(uint16) error
	Deallocate(uint16) error
}

// Acquire returns the default controller implementation
func Acquire() (Controller, error) {
	// TODO, make this a daemon, not filesystem implementation
	return filesystem.NewFSController("")
}
