package rdtp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestComputeSum(t *testing.T) {
	payload := []byte{0x61, 0x61, 0x61}
	p := &Packet{
		SrcPort: uint16(8080),
		DstPort: uint16(8081),
		Length:  uint16(len(payload)),
		Payload: payload,
	}
	assert.Equal(t, p.computeChecksum(), ^uint16(len(payload)+8080+8081+0x61*len(payload)))
}

func TestCheck(t *testing.T) {
	payload := []byte{0x61, 0x61, 0x61}
	p := &Packet{
		SrcPort:  uint16(8080),
		DstPort:  uint16(8081),
		Length:   uint16(len(payload)),
		Payload:  payload,
		Checksum: ^uint16(len(payload) + 8080 + 8081 + 0x61*len(payload)),
	}
	assert.True(t, p.Check())
	// tamper with the packet
	p.SrcPort = uint16(8000)
	assert.False(t, p.Check())
}

func TestComputeAndCheck(t *testing.T) {
	payload := []byte("[ mock http request ]")
	malformed := []byte("[ mock h0tp request ]")
	p := &Packet{
		SrcPort: uint16(8080),
		DstPort: uint16(8081),
		Length:  uint16(len(payload)),
		Payload: payload,
	}
	p.Checksum = p.computeChecksum()
	assert.True(t, p.Check())
	p.Payload = malformed
	assert.False(t, p.Check())
}
