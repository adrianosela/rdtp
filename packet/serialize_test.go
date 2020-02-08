package packet

import (
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSerialize(t *testing.T) {
	payload := []byte("[ mock http request ]")

	p, err := NewPacket(uint16(8081), uint16(8082), payload)
	assert.Nil(t, err)

	header := make([]byte, HeaderByteSize)
	binary.BigEndian.PutUint16(header[0:2], p.SrcPort)
	binary.BigEndian.PutUint16(header[2:4], p.DstPort)
	binary.BigEndian.PutUint16(header[4:6], p.Length)
	binary.BigEndian.PutUint16(header[6:8], p.Checksum)
	header[8] = uint8(0) // flags

	byt := p.Serialize()
	assert.Equal(t, string(byt), string(append(header, payload...)))
}

func TestDeserialize(t *testing.T) {

	payload := []byte("[ mock http request ]")

	serialized := append([]byte{
		31, 145, 31, 146, // src port, dst port
		0, byte(len(payload)), 185, 28, // length, checksum
		0}, payload...) // flags, payload

	p, err := Deserialize(serialized)
	assert.Nil(t, err)

	assert.Equal(t, p.SrcPort, uint16(8081))
	assert.Equal(t, p.DstPort, uint16(8082))
	assert.Equal(t, p.Length, uint16(len(payload)))
	assert.Equal(t, p.Checksum, p.computeChecksum())

	// ensure we dont deserialize non-packet data
	_, err = Deserialize([]byte("small"))
	assert.NotNil(t, err)

	// ensure we dont deserialize packets with bad size
	badLength := append([]byte{
		31, 145, 31, 146, // src port, dst port
		0, byte(len(payload)) + 1, 185, 28, // length, checksum
	}, payload...) // payload

	_, err = Deserialize(badLength)
	assert.NotNil(t, err)
}

func TestSerializeDeserialize(t *testing.T) {
	pLocal, err := NewPacket(uint16(8081), uint16(8082), []byte("[ mock http request ]"))
	assert.Nil(t, err)

	byt := pLocal.Serialize()

	pRemote, err := Deserialize(byt)
	assert.Nil(t, err)

	assert.EqualValues(t, pRemote, pLocal)
}
