package ports

// Manager is the interface in charge of managing
// ports for a given rdtp implementation
type Manager interface {
	AllocateAny() (uint16, error)
	Allocate(uint16) error
	Deallocate(uint16) error
}
