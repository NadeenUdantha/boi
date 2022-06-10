package boivm

type BoiMem struct {
	X []uint8
}

func (m *BoiMem) check(addr uint16) uint16 {
	return uint16(addr & 0x7fff)
}

func (m *BoiMem) r8(addr uint16) uint8 {
	addr = m.check(addr)
	return m.X[addr]
}

func (m *BoiMem) r16(addr uint16) uint16 {
	addr = m.check(addr)
	return uint16(m.X[addr]) | uint16(m.X[addr+1])<<8
}

func (m *BoiMem) av(x uint16) uint16 {
	if x&0x8000 == 0x8000 {
		return m.r16(x & 0x7fff)
	}
	return x
}

func (m *BoiMem) w8(addr uint16, v uint8) {
	addr = m.check(addr)
	m.X[addr] = v
}

func (m *BoiMem) w16(addr uint16, v uint16) {
	addr = m.check(addr)
	m.X[addr] = uint8(v)
	m.X[addr+1] = uint8(v << 8)
}
