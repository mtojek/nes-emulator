package main

type cpu6502 struct {
	a      uint8  // Accumulator Register
	x      uint8  // X Register
	y      uint8  // Y Register
	sp     uint8  // Stack Pointer (points to location on bus)
	pc     uint16 // Program Counter
	status uint8  // Status Register

	bus readableWriteable

	fetched uint8
	addrAbs uint16
	addrRel uint16
	opcode  uint8
	cycles  uint8

	// Lookups
	lookupOpcodes []instruction
}

func createCPU(b readableWriteable) *cpu6502 {
	c := &cpu6502{
		bus: b,
	}
	c.setupInstructions()
	return c
}

func (c *cpu6502) clock() {
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

func (c *cpu6502) read(addr uint16) uint8 {
	return c.bus.read(addr, true)
}

func (c *cpu6502) fetch() uint8 {
	if c.lookupOpcodes[c.opcode].addressing.name == lblAddressingModeIMP {
		c.fetched = c.read(c.addrAbs)
	}
	return c.fetched
}

func (c *cpu6502) reset() {
	c.addrAbs = 0xFFFC
	c.pc = uint16(c.read(c.addrAbs + 1)) << 8 | uint16(c.read(c.addrAbs))
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