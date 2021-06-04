package nes

import (
	"github.com/mtojek/nes-emulator/controller"
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

	player1 *controller.Controller
	player2 *controller.Controller
}

func Create() *NES {
	var cpuBus bus.Bus
	var ppuBus bus.Bus

	// Controllers
	player1 := controller.NewController()
	player2 := controller.NewController()

	// CPU
	aCPU := cpu.Create(&cpuBus)

	cpuBus.Connect(0x0000, 0x1FFF, memory.CreateMirroring(memory.CreateMemory(), 0x0000, 0x07FF))
	cpuBus.Connect(0x4016, 0x4016, player1)
	cpuBus.Connect(0x4017, 0x4017, player2)

	// PPU
	aPPU := ppu.Create(&cpuBus, &ppuBus, aCPU)

	cpuBusConnector := aPPU.CPUBusConnector()
	cpuBus.Connect(0x2000, 0x3FFF, memory.CreateMirroring(cpuBusConnector, 0x2000, 0x07))
	cpuBus.Connect(0x4014, 0x4014, cpuBusConnector) // DMA

	ppuBusConnector := aPPU.PPUBusConnector()
	ppuBus.Connect(0x2000, 0x3EFF, memory.CreateMirroring(ppuBusConnector, 0x2000, 0x0FFF))
	ppuBus.Connect(0x3F00, 0x3FFF, memory.CreateMirroring(ppuBusConnector, 0x3F00, 0x1F))

	return &NES{
		cpuBus: &cpuBus,
		ppuBus: &ppuBus,

		cpu: aCPU,
		ppu: aPPU,

		player1: player1,
		player2: player2,
	}
}

func (n *NES) Insert(cart *cartridge.Cartridge) {
	cart.ConnectTo(n.cpuBus, n.ppuBus)
	n.ppu.SetMirroring(cart)
}

func (n *NES) Reset() {
	n.cpu.Reset()
	n.ppu.Reset()
	n.systemClock = 0
}

func (n *NES) Clock() {
	n.ppu.Clock()

	if n.systemClock%3 == 0 {
		n.cpu.Clock()
	}

	if n.ppu.TriggerNMI {
		n.ppu.TriggerNMI = false
		n.cpu.NMI()
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

func (n *NES) UpdateControllers(state1, state2 [8]bool) {
	n.player1.SetState(state1)
	n.player2.SetState(state2)
}
