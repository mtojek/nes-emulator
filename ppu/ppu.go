package ppu

import (
	"image"
	"image/color"
	"math/rand"

	"github.com/mtojek/nes-emulator/bus"
)

type PPU2C02 struct {
	nametable [2][1024]uint8
	patterns  [2][4096]uint8
	palette   [32]uint8

	front *image.RGBA

	frameComplete bool

	scanline int16
	cycle    uint16

	controlReg   uint8bits
	maskReg      uint8bits
	statusReg    uint8bits
	addressLatch uint8
	dataBuffer   uint8bits

	vramAddrReg uint16 // Active "pointer" address into nametable to extract background tile info
	tramAddrReg uint16 // Temporary store of information to be "transferred" into "pointer" at various times

	cpuBus bus.ReadableWriteable
	ppuBus bus.ReadableWriteable
}

type cpuBusConnector struct {
	ppu *PPU2C02
}

func (cbc *cpuBusConnector) Read(addr uint16, bReadOnly bool) uint8 {
	//fmt.Printf("(implement me) read addr: %04x\n", addr)

	var data uint8bits
	if bReadOnly {
		switch addr {
		case 0x0000: // Control
			data = cbc.ppu.controlReg
		case 0x0001: // Mask
			data = cbc.ppu.maskReg
		case 0x0002: // Status
			data = cbc.ppu.statusReg
		case 0x0003: // OAM Address
		case 0x0004: // OAM Data
		case 0x0005: // Scroll
		case 0x0006: // PPU Address
		case 0x0007: // PPU Data
		}
	} else {
		switch addr {
		case 0x0000: // Control - Not readable
		case 0x0001: // Mask - Not Readable
		case 0x0002: // Status
			data = (cbc.ppu.statusReg & 0xE0) | (cbc.ppu.dataBuffer) & 0x1F
			cbc.ppu.statusReg = cbc.ppu.statusReg.withBit(flagStatusVerticalBlank, false)
			cbc.ppu.addressLatch = 0
		case 0x0003: // OAM Address
		case 0x0004: // OAM Data
		case 0x0005: // Scroll - Not Readable
		case 0x0006: // PPU Address - Not Readable
		case 0x0007: // PPU Data
			data = cbc.ppu.dataBuffer
			cbc.ppu.dataBuffer = uint8bits(cbc.ppu.ppuBus.Read(cbc.ppu.vramAddrReg, false))
			if cbc.ppu.vramAddrReg >= 0x3F00 {
				data = cbc.ppu.dataBuffer
			}

			if cbc.ppu.controlReg.bit(flagControlIncrementMode) {
				cbc.ppu.vramAddrReg += 32
			} else {
				cbc.ppu.vramAddrReg += 1
			}
		}
	}
	return uint8(data)
}

func (cbc *cpuBusConnector) Write(addr uint16, data uint8) {
	//panic("implement me")
	//fmt.Printf("(implement me) write addr: %04x, data: %02x\n", addr, data)
}

type ppuBusConnector struct {
	ppu *PPU2C02
}

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

	p.front.Set(int(p.cycle-1), int(p.scanline), nesPalette[colorOffest])

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
	return &cpuBusConnector{p}
}

func (p *PPU2C02) PPUBusConnector() bus.ReadableWriteable {
	return &ppuBusConnector{p}
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

func (p *PPU2C02) colourFromPaletteRAM(paletteIndex uint8, pixel uint8) color.RGBA {
	return nesPalette[p.ppuBus.Read(0x3F00+(uint16(paletteIndex)<<2)+uint16(pixel), false)&0x3F]
}
