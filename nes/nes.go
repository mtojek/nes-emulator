package nes

import (
	"image"

	"github.com/mtojek/nes-emulator/bus"
	"github.com/mtojek/nes-emulator/cartridge"
	"github.com/mtojek/nes-emulator/cpu"
	"github.com/mtojek/nes-emulator/memory"
	"github.com/mtojek/nes-emulator/ppu"
)

type NES struct {
	systemClock uint64

	cpuBus *bus.Bus
	ppuBus *bus.Bus

	cpu *cpu.CPU6502
	ppu *ppu.PPU2C02
}

func Create() *NES {
	var cpuBus bus.Bus
	var ppuBus bus.Bus

	// CPU
	aCPU := cpu.Create(&cpuBus)

	cpuBus.Connect(0x0000, 0x1FFF, memory.CreateMirroring(memory.CreateMemory(), 0x07FF))

	// PPU
	aPPU := ppu.Create(&cpuBus, &ppuBus)

	ppuBusConnector := aPPU.PPUBusConnector()
	ppuBus.Connect(0x0000, 0x1FFF, ppuBusConnector)
	ppuBus.Connect(0x2000, 0x3EFF, memory.CreateMirroring(ppuBusConnector, 0x0FFF))
	ppuBus.Connect(0x3F00, 0x3FFF, memory.CreateMirroring(ppuBusConnector, 0x0FFF))

	cpuBusConnector := aPPU.CPUBusConnector()
	cpuBus.Connect(0x2000, 0x3FFF, memory.CreateMirroring(cpuBusConnector, 0x001F))

	return &NES{
		cpuBus: &cpuBus,
		ppuBus: &ppuBus,

		cpu: aCPU,
		ppu: aPPU,
	}
}

func (n *NES) Insert(cart *cartridge.Cartridge) {
	cart.ConnectTo(n.cpuBus, n.ppuBus)
}

func (n *NES) Reset() {
	n.cpu.Reset()
	n.systemClock = 0
}

func (n *NES) Clock() {
	n.ppu.Clock()
	if n.systemClock%3 == 0 {
		n.cpu.Clock()
	}
	n.systemClock++
}

func (n *NES) FrameComplete() bool {
	return n.ppu.FrameComplete()
}

func (n *NES) DrawNewFrame() {
	n.ppu.DrawNewFrame()
}

func (n *NES) Buffer() *image.RGBA {
	return n.ppu.Buffer()
}
