package mappers

import "github.com/mtojek/nes-emulator/bus"

type mapper002 struct {
	nPRGBanks uint8
	prgBank1  uint8
	prgBank2  uint8
}

var _ Mapper = new(mapper002)

func (m mapper002) ID() uint8 {
	return 2
}

func (m mapper002) ConnectTo(cpuBus *bus.Bus, ppuBus *bus.Bus, prgMemory bus.ReadableWriteable, chrMemory bus.ReadableWriteable) {
	cpuBus.Connect(0x8000, 0xFFFF, prgMemory)
	ppuBus.Connect(0x0000, 0x1FFF, chrMemory)
}

func (m mapper002) MapCPU(addr uint16) uint64 {
	if addr >= 0xC000 {
		offset := uint64(m.prgBank2)*0x4000 + (uint64(addr) - 0xC000)
		return offset
	}
	if addr >= 0x8000 {
		return uint64(m.prgBank1)*0x4000 + (uint64(addr) - 0x8000)
	}
	panic("range not implemented")
}

func (m mapper002) MapPPU(addr uint16) uint64 {
	return uint64(addr)
}
