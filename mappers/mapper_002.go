package mappers

import "github.com/mtojek/nes-emulator/bus"

type mapper002 struct {
	nPRGBanks uint8
	prgBank1  uint8
	prgBank2  uint8
}

func (m *mapper002) CPURead(addr uint16) (uint64, bool) {
	if addr >= 0xC000 {
		offset := uint64(m.prgBank2)*0x4000 + (uint64(addr) - 0xC000)
		return offset, false
	}
	if addr >= 0x8000 {
		offset := uint64(m.prgBank1)*0x4000 + (uint64(addr) - 0x8000)
		return offset, false
	}
	panic("CPURead: range not implemented")
}

func (m *mapper002) CPUWrite(addr uint16, value uint8) (uint64, bool) {
	if addr > 0x8000 {
		m.prgBank1 = value % m.nPRGBanks
		return uint64(addr), true
	}
	return uint64(addr), false
}

func (m *mapper002) PPURead(addr uint16) (uint64, bool) {
	return uint64(addr), false
}

func (m *mapper002) PPUWrite(addr uint16, value uint8) (uint64, bool) {
	return m.PPURead(addr)
}

var _ Mapper = new(mapper002)

func (m mapper002) ID() uint8 {
	return 2
}

func (m *mapper002) ConnectTo(cpuBus *bus.Bus, ppuBus *bus.Bus, prgMemory bus.ReadableWriteable, chrMemory bus.ReadableWriteable) {
	cpuBus.Connect(0x8000, 0xFFFF, prgMemory)
	ppuBus.Connect(0x0000, 0x1FFF, chrMemory)
}
