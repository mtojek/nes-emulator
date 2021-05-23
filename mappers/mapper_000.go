package mappers

import "github.com/mtojek/nes-emulator/bus"

type mapper000 struct{}

var _ Mapper = new(mapper000)

func (m mapper000) ID() uint8 {
	return 0
}

func (m mapper000) ConnectTo(cpuBus *bus.Bus, ppuBus *bus.Bus, cart bus.ReadableWriteable) {
	cpuBus.Connect(0x8000, 0xFFFF, cart)
	ppuBus.Connect(0x0000, 0x1FFF, cart)
}

func (m mapper000) Map(addr uint16) uint16 {
	panic("TODO")
}
