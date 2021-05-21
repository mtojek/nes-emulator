package main

const (
	flagC = uint8(1) << 0 // Carry Bit
	flagZ = uint8(1) << 1 // Zero
	flagI = uint8(1) << 2 // Disable Interrupts
	flagD = uint8(1) << 3 // Decimal Mode (unused here)
	flagB = uint8(1) << 4 // Break
	flagU = uint8(1) << 5 // Unused
	flagV = uint8(1) << 6 // Overflow
	flagN = uint8(1) << 7 // Negative
)

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
	lookupAddressingModes map[string]addressingModeFunc
	lookupOpcodes         map[uint8]instruction
}

func createCPU(b readableWriteable) *cpu6502 {
	c := &cpu6502{
		bus: b,
	}
	c.lookupAddressingModes = map[string]addressingModeFunc{
		addressingModeIMP: c.imp, addressingModeIMM: c.imm, addressingModeZP0: c.zp0,
		addressingModeZPX: c.zpx, addressingModeZPY: c.zpy, addressingModeABX: c.abx,
		addressingModeABY: c.aby, addressingModeIND: c.ind, addressingModeIZX: c.izx,
		addressingModeIZY: c.izy, addressingModeREL: c.rel, addressingModeABS: c.abs,
	}
	c.lookupOpcodes = map[uint8]instruction{
		
	}
	return c
}

func (c *cpu6502) clock() {
	if c.cycles == 0 {
		c.opcode = c.read(c.pc)
		c.pc++

		ins := c.lookupOpcodes[c.opcode]

		c.cycles = ins.cycles
		c1 := ins.addressingModeDo()
		c2 := ins.operationDo()
		c.cycles += c1 & c2
	}
	c.cycles--
}

func (c *cpu6502) read(addr uint16) uint8 {
	return c.bus.read(addr, true)
}

func (c *cpu6502) setFlag(f uint8, value bool) {
	if value {
		c.status |= f
		return
	}
	c.status &= ^f
}
