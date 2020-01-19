package rdtp

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
	assert.Equal(t, p.Checksum, p.computeChecksum())
}

// Test the payload size limits are defined exactly
func TestLimits(t *testing.T) {
	_, err := NewPacket(uint16(8081), uint16(8082), make([]byte, MaxPayloadByteSize))
	assert.Nil(t, err)
	_, err = NewPacket(uint16(8081), uint16(8082), make([]byte, MaxPayloadByteSize+1))
	assert.NotNil(t, err)
}
