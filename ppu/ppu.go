package ppu

import (
	"github.com/mtojek/nes-emulator/bus"
	"image"
	"math/rand"
)

type PPU2C02 struct {
	nametable [2][1024]uint8
	patterns  [2][4096]uint8
	palette   [32]uint8

	front *image.RGBA

	frameComplete bool

	scanline int16
	cycle    uint16

	cpuBus bus.ReadableWriteable
	ppuBus bus.ReadableWriteable
}

type cpuBusConnector struct{}

func (cbc *cpuBusConnector) Read(addr uint16, bReadOnly bool) uint8 {
	//panic("implement me")
	//fmt.Printf("(implement me) read addr: %04x\n", addr)
	return 0
}

func (cbc *cpuBusConnector) Write(addr uint16, data uint8) {
	//panic("implement me")
	//fmt.Printf("(implement me) write addr: %04x, data: %02x\n", addr, data)
}

type ppuBusConnector struct{}

func (pbc *ppuBusConnector) Read(addr uint16, bReadOnly bool) uint8 {
	//panic("implement me")
	//fmt.Printf("(implement me) read addr: %04x\n", addr)
	return 0
}

func (pbc *ppuBusConnector) Write(addr uint16, data uint8) {
	//panic("implement me")
	//fmt.Printf("(implement me) write addr: %04x, data: %02x\n", addr, data)
}

var _ bus.ReadableWriteable = new(cpuBusConnector)

func Create(cpuBus, ppuBus bus.ReadableWriteable) *PPU2C02 {
	return &PPU2C02{
		cpuBus: cpuBus,
		ppuBus: ppuBus,

		front: image.NewRGBA(image.Rect(0, 0, 256, 240)),
	}
}

func (p *PPU2C02) Clock() {
	var colorOffest int
	if rand.Int()%2 == 1 {
		colorOffest = 0x3F
	} else {
		colorOffest = 0x30
	}

	p.front.Set(int(p.cycle-1), int(p.scanline), palette[colorOffest])

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

func (p *PPU2C02) CPUBusConnector() bus.ReadableWriteable {
	return new(cpuBusConnector)
}

func (p *PPU2C02) PPUBusConnector() bus.ReadableWriteable {
	return new(ppuBusConnector)
}

func (p *PPU2C02) DrawNewFrame() {
	p.frameComplete = false
}

func (p *PPU2C02) FrameComplete() bool {
	return p.frameComplete
}

func (p *PPU2C02) Buffer() *image.RGBA {
	return p.front
}
