package cpu

import (
	"github.com/mtojek/nes-emulator/bus"
)

type CPU6502 struct {
	a      uint8  // Accumulator Register
	x      uint8  // X Register
	y      uint8  // Y Register
	sp     uint8  // Stack Pointer (points to location on bus)
	pc     uint16 // Program Counter
	status uint8  // Status Register

	fetched uint8
	addrAbs uint16
	addrRel uint16
	opcode  uint8
	cycles  uint8

	// Lookups
	lookupOpcodes []instruction

	bus bus.ReadableWriteable
}

func Create(b bus.ReadableWriteable) *CPU6502 {
	c := &CPU6502{
		bus: b,
	}
	c.setupInstructions()
	c.Reset()
	return c
}

func (c *CPU6502) Reset() {
	c.addrAbs = 0xFFFC
	c.pc = uint16(c.read(c.addrAbs+1))<<8 | uint16(c.read(c.addrAbs))
	c.a = 0
	c.x = 0
	c.y = 0
	c.sp = 0xFD
	c.status = 0 | flagU
	c.addrRel = 0
	c.addrAbs = 0
	c.fetched = 0
	c.cycles = 8
}

func (c *CPU6502) Clock() {
	if c.cycles == 0 {
		c.opcode = c.read(c.pc)
		c.pc++

		ins := c.lookupOpcodes[c.opcode]
		c.cycles = ins.cycles
		c1 := ins.addressing.do()
		c2 := ins.operationDo()
		c.cycles += c1 & c2
	}
	c.cycles--
}

func (c *CPU6502) read(addr uint16) uint8 {
	return c.bus.Read(addr, true)
}

func (c *CPU6502) write(addr uint16, data uint8) {
	c.bus.Write(addr, data)
}

func (c *CPU6502) fetch() uint8 {
	ins := c.lookupOpcodes[c.opcode]
	if ins.addressing.name != lblAddressingModeIMP {
		c.fetched = c.read(c.addrAbs)
	}
	return c.fetched
}
