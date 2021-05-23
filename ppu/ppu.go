package ppu

import "github.com/mtojek/nes-emulator/bus"

type PPU2C02 struct {
	nametable [2][1024]uint8
	palette   [32]uint8
	patterns  [2][4096]uint8

	cpuBus bus.ReadableWriteable
	ppuBus bus.ReadableWriteable
}

func Create(cpuBus, ppuBus bus.ReadableWriteable) *PPU2C02 {
	return &PPU2C02{
		cpuBus: cpuBus,
		ppuBus: ppuBus,
	}
}

func (p *PPU2C02) Registers() bus.ReadableWriteable {
	panic("TODO")
}
