package packet

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetSum(t *testing.T) {
	payload := []byte{0x61, 0x61, 0x61}
	p := &Packet{
		SrcPort: uint16(8080),
		DstPort: uint16(8081),
		SeqNo:   uint32(10),
		AckNo:   uint32(11),
		Length:  uint16(len(payload)),
		Payload: payload,
	}
	p.SetSum()
	assert.Equal(t, p.Checksum, ^uint16(len(payload)+8080+8081+10+11+0x61*len(payload)))
}

func TestCheckSum(t *testing.T) {
	payload := []byte{0x61, 0x61, 0x61}
	p := &Packet{
		SrcPort:  uint16(8080),
		DstPort:  uint16(8081),
		SeqNo:    uint32(10),
		AckNo:    uint32(11),
		Length:   uint16(len(payload)),
		Payload:  payload,
		Checksum: ^uint16(len(payload) + 8080 + 8081 + 10 + 11 + 0x61*len(payload)),
	}

	assert.True(t, p.CheckSum())
	// tamper with the packet
	p.SrcPort = uint16(8000)
	assert.False(t, p.CheckSum())
}

func TestSum(t *testing.T) {
	payload := []byte("[ mock http request ]")
	malformed := []byte("[ mock h0tp request ]")
	p := &Packet{
		SrcPort: uint16(8080),
		DstPort: uint16(8081),
		SeqNo:   uint32(10),
		AckNo:   uint32(11),
		Length:  uint16(len(payload)),
		Payload: payload,
	}
	p.Checksum = p.sum()
	assert.True(t, p.CheckSum())
	p.Payload = malformed
	assert.False(t, p.CheckSum())
}
