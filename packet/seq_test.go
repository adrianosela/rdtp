package packet

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetSeqNo(t *testing.T) {
	p, err := NewPacket(uint16(14), uint16(15), nil)
	assert.Nil(t, err)

	testSeq := uint32(28364)

	p.SetSeqNo(testSeq)
	assert.Equal(t, testSeq, p.SeqNo)
}

func TestSetAckNo(t *testing.T) {
	p, err := NewPacket(uint16(14), uint16(15), nil)
	assert.Nil(t, err)

	testAck := uint32(28364)

	p.SetAckNo(testAck)
	assert.Equal(t, testAck, p.AckNo)
}
