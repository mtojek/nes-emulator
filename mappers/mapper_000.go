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

func (m mapper000) MapCPU(addr uint16) uint64 {
	if m.nPRGBanks > 1 {
		return uint64(addr & 0x7FFF)
	}
	return uint64(addr & 0x3FFF)
}

func (m mapper000) MapPPU(addr uint16) uint64 {
	return uint64(addr)
}
