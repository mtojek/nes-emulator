package ppu

import (
	"github.com/mtojek/nes-emulator/bus"
)

type PPU2C02 struct {
	nametable [2][1024]uint8
	patterns  [2][4096]uint8
	palette   [32]uint8

	frameComplete bool

	scanline int16
	cycle uint16

	cpuBus bus.ReadableWriteable
	ppuBus bus.ReadableWriteable
}

type registersHandler struct{}


func (rh *registersHandler) Read(addr uint16, bReadOnly bool) uint8 {
	//panic("implement me")
	//fmt.Printf("(implement me) read addr: %04x\n", addr)
	return 0
}

func (rh *registersHandler) Write(addr uint16, data uint8) {
	//panic("implement me")
	//fmt.Printf("(implement me) write addr: %04x, data: %02x\n", addr, data)
}

var _ bus.ReadableWriteable = new(registersHandler)

func Create(cpuBus, ppuBus bus.ReadableWriteable) *PPU2C02 {
	return &PPU2C02{
		cpuBus: cpuBus,
		ppuBus: ppuBus,
	}
}

func (p *PPU2C02) Clock() {
	// fake noise
	// sprScreen.SetPixel(cycle - 1, scanline, palScreen[(rand() % 2) ? 0x3F : 0x30]);

	p.cycle++
	if p.cycle >= 341 {
		p.cycle = 0
		p.scanline++
		if p.scanline >= 261 {
			p.scanline = -1
			p.frameComplete = true
		}
	}
}

func (p *PPU2C02) Registers() bus.ReadableWriteable {
	return new(registersHandler)
}

func (p *PPU2C02) DrawNewFrame() {
	p.frameComplete = false
}

func (p *PPU2C02) FrameComplete() bool {
	return p.frameComplete
}