package packet

import (
	"testing"

	"github.com/google/gopacket/layers"
	"github.com/stretchr/testify/assert"
)

func TestIPv4Details(t *testing.T) {
	p, err := NewPacket(uint16(14), uint16(15), nil)
	assert.Nil(t, err)

	// get error when not set
	_, err = p.IPv4Details()
	assert.NotNil(t, err)

	// set and then get no error
	ipv4Set := &layers.IPv4{}
	p.SetIPv4Details(ipv4Set)

	ipv4Got, err := p.IPv4Details()
	assert.Nil(t, err)

	assert.EqualValues(t, ipv4Set, ipv4Got)
}
