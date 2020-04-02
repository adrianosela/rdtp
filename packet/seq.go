package packet

// SetSeqNo sets the sequence number of this packet
func (p *Packet) SetSeqNo(seq uint32) {
	p.SeqNo = seq
}

// SetAckNo sets the acknowledgement number of this packet
func (p *Packet) SetAckNo(ack uint32) {
	p.AckNo = ack
}
