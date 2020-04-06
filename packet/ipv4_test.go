package packet

import (
	"testing"

	"github.com/google/gopacket/layers"
	"github.com/stretchr/testify/assert"
)

func TestIPv4(t *testing.T) {
	p, err := NewPacket(uint16(14), uint16(15), nil)
	assert.Nil(t, err)

	// get error when not set
	_, err = p.IPv4()
	assert.NotNil(t, err)

	// set and then get no error
	ipv4Set := &layers.IPv4{}
	p.SetIPv4(ipv4Set)

	ipv4Got, err := p.IPv4()
	assert.Nil(t, err)

	assert.EqualValues(t, ipv4Set, ipv4Got)
}
