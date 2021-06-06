package mappers

import "github.com/mtojek/nes-emulator/bus"

type mapper000 struct {
	nPRGBanks uint8
}

var _ Mapper = new(mapper000)

func (m mapper000) ID() uint8 {
	return 0
}

func (m mapper000) ConnectTo(cpuBus *bus.Bus, ppuBus *bus.Bus, prgMemory bus.ReadableWriteable, chrMemory bus.ReadableWriteable) {
	cpuBus.Connect(0x8000, 0xFFFF, prgMemory)
	ppuBus.Connect(0x0000, 0x1FFF, chrMemory)
}

func (m *mapper000) CPURead(addr uint16) (uint64, bool) {
	if m.nPRGBanks > 1 {
		return uint64(addr & 0x7FFF), false
	}
	return uint64(addr & 0x3FFF), false
}

func (m *mapper000) CPUWrite(addr uint16, value uint8) (uint64, bool) {
	return m.CPURead(addr)
}

func (m *mapper000) PPURead(addr uint16) (uint64, bool) {
	return uint64(addr), false
}

func (m *mapper000) PPUWrite(addr uint16, value uint8) (uint64, bool) {
	return m.PPURead(addr)
}
