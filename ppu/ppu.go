package ppu

import (
	"github.com/mtojek/nes-emulator/bus"
	"image"
	"image/color"
)

const (
	HORIZONTAL   = 0
	VERTICAL     = 1
	ONESCREEN_LO = 2
	ONESCREEN_HI = 3
)

type PPU2C02 struct {
	cpuBus bus.ReadableWriteable
	ppuBus bus.ReadableWriteable

	nametable [2][1024]uint8
	patterns  [2][4096]uint8
	palette   [32]uint8

	background *image.RGBA

	frameComplete bool

	scanline int16
	cycle    uint16

	// PPU registers
	v uint16 // current vram address (15 bit)
	t uint16 // temporary vram address (15 bit)
	x byte   // fine x scroll (3 bit)
	w byte   // write toggle (1 bit)
	f byte   // even/odd frame flag (1 bit)

	register byte

	// background temporary variables
	nameTableByte      byte
	attributeTableByte byte
	lowTileByte        byte
	highTileByte       byte
	tileData           uint64

	// $2000 PPUCTRL
	flagNameTable       byte // 0: $2000; 1: $2400; 2: $2800; 3: $2C00
	flagIncrement       byte // 0: add 1; 1: add 32
	flagSpriteTable     byte // 0: $0000; 1: $1000; ignored in 8x16 mode
	flagBackgroundTable byte // 0: $0000; 1: $1000
	flagSpriteSize      byte // 0: 8x8; 1: 8x16
	flagMasterSlave     byte // 0: read EXT; 1: write EXT

	// $2001 PPUMASK
	flagGrayscale          byte // 0: color; 1: grayscale
	flagShowLeftBackground byte // 0: hide; 1: show
	flagShowLeftSprites    byte // 0: hide; 1: show
	flagShowBackground     byte // 0: hide; 1: show
	flagShowSprites        byte // 0: hide; 1: show
	flagRedTint            byte // 0: normal; 1: emphasized
	flagGreenTint          byte // 0: normal; 1: emphasized
	flagBlueTint           byte // 0: normal; 1: emphasized

	// $2002 PPUSTATUS
	flagSpriteZeroHit  byte
	flagSpriteOverflow byte

	// $2003 OAMADDR
	oamAddress byte

	// $2007 PPUDATA
	bufferedData byte // for buffered reads

	// Mirroring
	mirroring uint8

	// Background
	bgNextTileId       uint8
	bgNextTileAttrib   uint8
	bgNextTileLsb      uint8
	bgNextTileMsb      uint8
	bgShifterPatternLo uint16
	bgShifterPatternHi uint16
	bgShifterAttribLo  uint16
	bgShifterAttribHi  uint16

	// NMI flags
	nmiOccurred bool
	nmiOutput   bool
	nmiPrevious bool
	nmiDelay    byte

	TriggerNMI bool
}

type mirrorer interface {
	Mirroring() uint8
}

type cpuBusConnector struct {
	ppu *PPU2C02
}

func (cbc *cpuBusConnector) Read(addr uint16) uint8 {
	//fmt.Printf("(implement me) read addr: %04x\n", addr)

	var data uint8
	switch addr {
	case 0x2000: // Control - Not readable
	case 0x2001: // Mask - Not Readable
	case 0x2002: // Status
		return cbc.ppu.readStatus()
	case 0x2003: // OAM Address
	case 0x2004: // OAM Data
	case 0x2005: // Scroll - Not Readable
	case 0x2006: // PPU Address - Not Readable
	case 0x2007: // PPU Data
		return cbc.ppu.readData()
	}
	return data
}

func (cbc *cpuBusConnector) Write(addr uint16, data uint8) {
	//fmt.Printf("(implement me) write addr: %04x, data: %02x\n", addr, data)
	cbc.ppu.register = data
	switch addr {
	case 0x2000: // Control
		cbc.ppu.writeControl(data)
	case 0x2001: // Mask
		cbc.ppu.writeMask(data)
	case 0x2002: // Status
	case 0x2003: // OAM Address
	case 0x2004: // OAM Data
	case 0x2005: // Scroll
		cbc.ppu.writeScroll(data)
	case 0x2006: // PPU Address
		cbc.ppu.writeAddress(data)
	case 0x2007: // PPU Data
		cbc.ppu.writeData(data)
	}
}

type ppuBusConnector struct {
	ppu *PPU2C02
}

func (pbc *ppuBusConnector) Read(addr uint16) uint8 {
	//fmt.Printf("(implement me) read addr: %04x\n", addr)

	var data uint8
	addr = addr & 0x3FFF

	if addr <= 0x1FFF {
		// If the cartridge cant map the address, have
		// a physical location ready here
		data = pbc.ppu.patterns[(addr&0x1000)>>12][addr&0x0FFF]
	} else if addr >= 0x2000 && addr <= 0x3EFF {
		addr &= 0x0FFF

		if pbc.ppu.mirroring == VERTICAL {
			// Vertical
			if addr <= 0x03FF {
				data = pbc.ppu.nametable[0][addr&0x03FF]
			} else if addr >= 0x0400 && addr <= 0x07FF {
				data = pbc.ppu.nametable[1][addr&0x03FF]
			} else if addr >= 0x0800 && addr <= 0x0BFF {
				data = pbc.ppu.nametable[0][addr&0x03FF]
			} else if addr >= 0x0C00 && addr <= 0x0FFF {
				data = pbc.ppu.nametable[1][addr&0x03FF]
			}
		} else if pbc.ppu.mirroring == HORIZONTAL {
			// Horizontal
			if addr <= 0x03FF {
				data = pbc.ppu.nametable[0][addr&0x03FF]
			} else if addr >= 0x0400 && addr <= 0x07FF {
				data = pbc.ppu.nametable[0][addr&0x03FF]
			} else if addr >= 0x0800 && addr <= 0x0BFF {
				data = pbc.ppu.nametable[1][addr&0x03FF]
			} else if addr >= 0x0C00 && addr <= 0x0FFF {
				data = pbc.ppu.nametable[1][addr&0x03FF]
			}
		}
	} else if addr >= 0x3F00 && addr <= 0x3FFF {
		addr &= 0x001F
		if addr == 0x0010 {
			addr = 0x0000
		}
		if addr == 0x0014 {
			addr = 0x0004
		}
		if addr == 0x0018 {
			addr = 0x0008
		}
		if addr == 0x001C {
			addr = 0x000C
		}

		c := uint8(0x3F)
		if pbc.ppu.flagGrayscale == 1 {
			c = 0x30
		}
		data = pbc.ppu.palette[addr] & c
	}
	return data
}

func (pbc *ppuBusConnector) Write(addr uint16, data uint8) {
	//fmt.Printf("(implement me) write addr: %04x, data: %02x\n", addr, data)

	addr &= 0x3FFF
	if addr <= 0x1FFF {
		pbc.ppu.patterns[(addr&0x1000)>>12][addr&0x0FFF] = data
	} else if addr >= 0x2000 && addr <= 0x3EFF {
		addr &= 0x0FFF
		if pbc.ppu.mirroring == VERTICAL {
			// Vertical
			if addr <= 0x03FF {
				pbc.ppu.nametable[0][addr&0x03FF] = data
			} else if addr >= 0x0400 && addr <= 0x07FF {
				pbc.ppu.nametable[1][addr&0x03FF] = data
			} else if addr >= 0x0800 && addr <= 0x0BFF {
				pbc.ppu.nametable[0][addr&0x03FF] = data
			} else if addr >= 0x0C00 && addr <= 0x0FFF {
				pbc.ppu.nametable[1][addr&0x03FF] = data
			}
		} else if pbc.ppu.mirroring == HORIZONTAL {
			// Horizontal
			if addr <= 0x03FF {
				pbc.ppu.nametable[0][addr&0x03FF] = data
			} else if addr >= 0x0400 && addr <= 0x07FF {
				pbc.ppu.nametable[0][addr&0x03FF] = data
			} else if addr >= 0x0800 && addr <= 0x0BFF {
				pbc.ppu.nametable[1][addr&0x03FF] = data
			} else if addr >= 0x0C00 && addr <= 0x0FFF {
				pbc.ppu.nametable[1][addr&0x03FF] = data
			}
		}
	} else if addr >= 0x3F00 && addr <= 0x3FFF {
		addr &= 0x001F
		if addr == 0x0010 {
			addr = 0x0000
		}
		if addr == 0x0014 {
			addr = 0x0004
		}
		if addr == 0x0018 {
			addr = 0x0008
		}
		if addr == 0x001C {
			addr = 0x000C
		}
		pbc.ppu.palette[addr] = data
	}
}

var _ bus.ReadableWriteable = new(cpuBusConnector)

func Create(cpuBus, ppuBus bus.ReadableWriteable) *PPU2C02 {
	return &PPU2C02{
		cpuBus: cpuBus,
		ppuBus: ppuBus,

		background: image.NewRGBA(image.Rect(0, 0, 256, 240)),
	}
}

func (ppu *PPU2C02) Reset() {
	ppu.cycle = 340
	ppu.scanline = 240
	ppu.writeControl(0)
	ppu.writeMask(0)
	//ppu.writeOAMAddress(0)
}

func (ppu *PPU2C02) tick() {
	if ppu.nmiDelay > 0 {
		ppu.nmiDelay--
		if ppu.nmiDelay == 0 && ppu.nmiOutput && ppu.nmiOccurred {
			ppu.TriggerNMI = true
		}
	}

	if ppu.flagShowBackground != 0 || ppu.flagShowSprites != 0 {
		if ppu.f == 1 && ppu.scanline == 261 && ppu.cycle == 339 {
			ppu.cycle = 0
			ppu.scanline = 0
			ppu.f ^= 1
			return
		}
	}

	ppu.cycle++
	if ppu.cycle > 340 {
		ppu.cycle = 0
		ppu.scanline++
		if ppu.scanline > 261 {
			ppu.scanline = 0
			ppu.f ^= 1
		}
	}
}

func (ppu *PPU2C02) renderPixel() {
	x := int(ppu.cycle - 1)
	y := int(ppu.scanline)
	background := ppu.backgroundPixel()

	if x < 8 && ppu.flagShowLeftBackground == 0 {
		background = 0
	}

	c := nesPalette[ppu.readPalette(uint16(background))%64]
	ppu.background.SetRGBA(x, y, c)
}

func (ppu *PPU2C02) backgroundPixel() byte {
	if ppu.flagShowBackground == 0 {
		return 0
	}
	data := ppu.fetchTileData() >> ((7 - ppu.x) * 4)
	return byte(data & 0x0F)
}

func (ppu *PPU2C02) fetchTileData() uint32 {
	return uint32(ppu.tileData >> 32)
}

func (ppu *PPU2C02) Clock() {
	ppu.tick()

	renderingEnabled := ppu.flagShowBackground != 0 || ppu.flagShowSprites != 0
	preLine := ppu.scanline == 261
	visibleLine := ppu.scanline < 240

	renderLine := preLine || visibleLine
	preFetchCycle := ppu.cycle >= 321 && ppu.cycle <= 336
	visibleCycle := ppu.cycle >= 1 && ppu.cycle <= 256
	fetchCycle := preFetchCycle || visibleCycle

	// background logic
	if renderingEnabled {
		if visibleLine && visibleCycle {
			ppu.renderPixel()
		}
		if renderLine && fetchCycle {
			ppu.tileData <<= 4
			switch ppu.cycle % 8 {
			case 1:
				ppu.fetchNameTableByte()
			case 3:
				ppu.fetchAttributeTableByte()
			case 5:
				ppu.fetchLowTileByte()
			case 7:
				ppu.fetchHighTileByte()
			case 0:
				ppu.storeTileData()
			}
		}
		if preLine && ppu.cycle >= 280 && ppu.cycle <= 304 {
			ppu.transferAddressY()
		}
		if renderLine {
			if fetchCycle && ppu.cycle%8 == 0 {
				ppu.incrementScrollX()
			}
			if ppu.cycle == 256 {
				ppu.incrementScrollY()
			}
			if ppu.cycle == 257 {
				ppu.transferAddressX()
			}
		}
	}

	// TODO sprite logic

	// vblank logic
	if ppu.scanline == 241 && ppu.cycle == 1 {
		ppu.setVerticalBlank()
	}
	if preLine && ppu.cycle == 1 {
		ppu.clearVerticalBlank()
		ppu.flagSpriteZeroHit = 0
		ppu.flagSpriteOverflow = 0
	}

	if preLine && ppu.cycle == 340 {
		ppu.frameComplete = true
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
	return p.background
}

func (p *PPU2C02) SetMirroring(m mirrorer) {
	p.mirroring = m.Mirroring()
}

func (p *PPU2C02) colourFromPaletteRAM(paletteIndex uint8, pixel uint8) color.RGBA {
	return nesPalette[p.ppuBus.Read(0x3F00+(uint16(paletteIndex)<<2)+uint16(pixel))&0x3F]
}

// $2002: PPUSTATUS
func (ppu *PPU2C02) readStatus() byte {
	result := ppu.register & 0x1F
	result |= ppu.flagSpriteOverflow << 5
	result |= ppu.flagSpriteZeroHit << 6
	if ppu.nmiOccurred {
		result |= 1 << 7
	}
	ppu.nmiOccurred = false
	ppu.nmiChange()
	// w:                   = 0
	ppu.w = 0
	return result
}

// $2007: PPUDATA (read)
func (ppu *PPU2C02) readData() byte {
	value := ppu.ppuBus.Read(ppu.v)
	// emulate buffered reads
	if ppu.v%0x4000 < 0x3F00 {
		buffered := ppu.bufferedData
		ppu.bufferedData = value
		value = buffered
	} else {
		ppu.bufferedData = ppu.ppuBus.Read(ppu.v - 0x1000)
	}
	// increment address
	if ppu.flagIncrement == 0 {
		ppu.v += 1
	} else {
		ppu.v += 32
	}
	return value
}

// $2000: PPUCTRL
func (ppu *PPU2C02) writeControl(value byte) {
	ppu.flagNameTable = (value >> 0) & 3
	ppu.flagIncrement = (value >> 2) & 1
	ppu.flagSpriteTable = (value >> 3) & 1
	ppu.flagBackgroundTable = (value >> 4) & 1
	ppu.flagSpriteSize = (value >> 5) & 1
	ppu.flagMasterSlave = (value >> 6) & 1
	ppu.nmiOutput = (value>>7)&1 == 1
	ppu.nmiChange()
	// t: ....BA.. ........ = d: ......BA
	ppu.t = (ppu.t & 0xF3FF) | ((uint16(value) & 0x03) << 10)
}

// $2001: PPUMASK
func (ppu *PPU2C02) writeMask(value byte) {
	ppu.flagGrayscale = (value >> 0) & 1
	ppu.flagShowLeftBackground = (value >> 1) & 1
	ppu.flagShowLeftSprites = (value >> 2) & 1
	ppu.flagShowBackground = (value >> 3) & 1
	ppu.flagShowSprites = (value >> 4) & 1
	ppu.flagRedTint = (value >> 5) & 1
	ppu.flagGreenTint = (value >> 6) & 1
	ppu.flagBlueTint = (value >> 7) & 1
}

// $2005: PPUSCROLL
func (ppu *PPU2C02) writeScroll(value byte) {
	if ppu.w == 0 {
		// t: ........ ...HGFED = d: HGFED...
		// x:               CBA = d: .....CBA
		// w:                   = 1
		ppu.t = (ppu.t & 0xFFE0) | (uint16(value) >> 3)
		ppu.x = value & 0x07
		ppu.w = 1
	} else {
		// t: .CBA..HG FED..... = d: HGFEDCBA
		// w:                   = 0
		ppu.t = (ppu.t & 0x8FFF) | ((uint16(value) & 0x07) << 12)
		ppu.t = (ppu.t & 0xFC1F) | ((uint16(value) & 0xF8) << 2)
		ppu.w = 0
	}
}

// $2006: PPUADDR
func (ppu *PPU2C02) writeAddress(value byte) {
	if ppu.w == 0 {
		// t: ..FEDCBA ........ = d: ..FEDCBA
		// t: .X...... ........ = 0
		// w:                   = 1
		ppu.t = (ppu.t & 0x80FF) | ((uint16(value) & 0x3F) << 8)
		ppu.w = 1
	} else {
		// t: ........ HGFEDCBA = d: HGFEDCBA
		// v                    = t
		// w:                   = 0
		ppu.t = (ppu.t & 0xFF00) | uint16(value)
		ppu.v = ppu.t
		ppu.w = 0
	}
}

// $2007: PPUDATA (write)
func (ppu *PPU2C02) writeData(value byte) {
	ppu.ppuBus.Write(ppu.v%0x4000, value)
	if ppu.flagIncrement == 0 {
		ppu.v += 1
	} else {
		ppu.v += 32
	}
}

// State machine
func (ppu *PPU2C02) incrementScrollX() {
	// increment hori(v)
	// if coarse X == 31
	if ppu.v&0x001F == 31 {
		// coarse X = 0
		ppu.v &= 0xFFE0
		// switch horizontal nametable
		ppu.v ^= 0x0400
	} else {
		// increment coarse X
		ppu.v++
	}
}

func (ppu *PPU2C02) incrementScrollY() {
	// increment vert(v)
	// if fine Y < 7
	if ppu.v&0x7000 != 0x7000 {
		// increment fine Y
		ppu.v += 0x1000
	} else {
		// fine Y = 0
		ppu.v &= 0x8FFF
		// let y = coarse Y
		y := (ppu.v & 0x03E0) >> 5
		if y == 29 {
			// coarse Y = 0
			y = 0
			// switch vertical nametable
			ppu.v ^= 0x0800
		} else if y == 31 {
			// coarse Y = 0, nametable not switched
			y = 0
		} else {
			// increment coarse Y
			y++
		}
		// put coarse Y back into v
		ppu.v = (ppu.v & 0xFC1F) | (y << 5)
	}
}

func (ppu *PPU2C02) transferAddressX() {
	// hori(v) = hori(t)
	// v: .....F.. ...EDCBA = t: .....F.. ...EDCBA
	ppu.v = (ppu.v & 0xFBE0) | (ppu.t & 0x041F)
}

func (ppu *PPU2C02) transferAddressY() {
	// vert(v) = vert(t)
	// v: .IHGF.ED CBA..... = t: .IHGF.ED CBA.....
	ppu.v = (ppu.v & 0x841F) | (ppu.t & 0x7BE0)
}

func (ppu *PPU2C02) setVerticalBlank() {
	//ppu.front, ppu.back = ppu.back, ppu.front
	ppu.nmiOccurred = true
	ppu.nmiChange()
}

func (ppu *PPU2C02) clearVerticalBlank() {
	ppu.nmiOccurred = false
	ppu.nmiChange()
}

func (ppu *PPU2C02) nmiChange() {
	nmi := ppu.nmiOutput && ppu.nmiOccurred
	if nmi && !ppu.nmiPrevious {
		// TODO: this fixes some games but the delay shouldn't have to be so
		// long, so the timings are off somewhere
		ppu.nmiDelay = 15
	}
	ppu.nmiPrevious = nmi
}

func (ppu *PPU2C02) fetchNameTableByte() {
	v := ppu.v
	address := 0x2000 | (v & 0x0FFF)
	ppu.nameTableByte = ppu.ppuBus.Read(address)
}

func (ppu *PPU2C02) fetchAttributeTableByte() {
	v := ppu.v
	address := 0x23C0 | (v & 0x0C00) | ((v >> 4) & 0x38) | ((v >> 2) & 0x07)
	shift := ((v >> 4) & 4) | (v & 2)
	ppu.attributeTableByte = ((ppu.ppuBus.Read(address) >> shift) & 3) << 2
}

func (ppu *PPU2C02) fetchLowTileByte() {
	fineY := (ppu.v >> 12) & 7
	table := ppu.flagBackgroundTable
	tile := ppu.nameTableByte
	address := 0x1000*uint16(table) + uint16(tile)*16 + fineY
	ppu.lowTileByte = ppu.ppuBus.Read(address)
}

func (ppu *PPU2C02) fetchHighTileByte() {
	fineY := (ppu.v >> 12) & 7
	table := ppu.flagBackgroundTable
	tile := ppu.nameTableByte
	address := 0x1000*uint16(table) + uint16(tile)*16 + fineY
	ppu.highTileByte = ppu.ppuBus.Read(address + 8)
}

func (ppu *PPU2C02) storeTileData() {
	var data uint32
	for i := 0; i < 8; i++ {
		a := ppu.attributeTableByte
		p1 := (ppu.lowTileByte & 0x80) >> 7
		p2 := (ppu.highTileByte & 0x80) >> 6
		ppu.lowTileByte <<= 1
		ppu.highTileByte <<= 1
		data <<= 4
		data |= uint32(a | p1 | p2)
	}
	ppu.tileData |= uint64(data)
}

func (ppu *PPU2C02) readPalette(address uint16) byte {
	if address >= 16 && address%4 == 0 {
		address -= 16
	}
	return ppu.palette[address]
}
