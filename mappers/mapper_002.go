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

func (m mapper002) MapCPU(addr uint16) uint16 {
	if addr >= 0xC000 {
		offset := int(m.prgBank2)*0x4000 + (int(addr) - 0xC000) // FIXME This must be able to return data from larger memory banks
		return uint16(offset)
	}
	if addr >= 0x8000 {
		return uint16(m.prgBank1)*0x4000 + (addr - 0x8000)
	}
	panic("range not implemented")
}

func (m mapper002) MapPPU(addr uint16) uint16 {
	return addr
}
