package packet

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFlags(t *testing.T) {
	p, err := NewPacket(uint16(14), uint16(15), nil)
	assert.Nil(t, err)

	tests := []struct {
		FlagName  string
		SetFunc   func()
		CheckFunc func() bool
	}{
		{
			FlagName:  "SYN",
			SetFunc:   func() { p.SetSYN() },
			CheckFunc: func() bool { return p.IsSYN() },
		},
		{
			FlagName:  "ACK",
			SetFunc:   func() { p.SetACK() },
			CheckFunc: func() bool { return p.IsACK() },
		},
		{
			FlagName:  "FIN",
			SetFunc:   func() { p.SetFIN() },
			CheckFunc: func() bool { return p.IsFIN() },
		},
		{
			FlagName:  "ERR",
			SetFunc:   func() { p.SetERR() },
			CheckFunc: func() bool { return p.IsERR() },
		},
	}

	for _, test := range tests {
		assert.False(t, test.CheckFunc())
		test.SetFunc()
		assert.True(t, test.CheckFunc())
	}
}
