package packet

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPacket(t *testing.T) {
	payload := []byte("[ mock http request ]")
	p, err := NewPacket(uint16(8081), uint16(8082), payload)
	assert.Nil(t, err)

	assert.Equal(t, p.SrcPort, uint16(8081))
	assert.Equal(t, p.DstPort, uint16(8082))
	assert.Equal(t, p.Length, uint16(len(payload)))
}

// Test the payload size limits are defined exactly
func TestLimits(t *testing.T) {
	_, err := NewPacket(uint16(8081), uint16(8082), make([]byte, MaxPayloadBytes))
	assert.Nil(t, err)
	_, err = NewPacket(uint16(8081), uint16(8082), make([]byte, MaxPayloadBytes+1))
	assert.NotNil(t, err)
}
