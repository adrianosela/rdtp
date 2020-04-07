package packet

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetDestinationIPv4(t *testing.T) {
	testIP := net.IP{127, 0, 0, 1}

	p, err := NewPacket(uint16(14), uint16(15), nil)
	assert.Nil(t, err)

	// error when not set
	_, err = p.GetDestinationIPv4()
	assert.NotNil(t, err)

	p.SetDestinationIPv4(testIP)
	got, err := p.GetDestinationIPv4()
	assert.Nil(t, err)
	assert.EqualValues(t, got, testIP)
}

func TestSetSourceIPv4(t *testing.T) {
	testIP := net.IP{127, 0, 0, 1}

	p, err := NewPacket(uint16(14), uint16(15), nil)
	assert.Nil(t, err)

	// error when not set
	_, err = p.GetSourceIPv4()
	assert.NotNil(t, err)

	p.SetSourceIPv4(testIP)
	got, err := p.GetSourceIPv4()
	assert.Nil(t, err)
	assert.EqualValues(t, got, testIP)
}
