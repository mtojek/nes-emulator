package main

// Implied addressing mode
func (c *cpu6502) imp() uint8 {
	c.fetched = c.a
	return 0
}

// Immediate addressing mode
func (c *cpu6502) imm() uint8 {
	c.addrAbs = c.pc
	c.pc++
	return 0
}

// Zero-page addressing
func (c *cpu6502) zp0() uint8 {
	c.addrAbs = uint16(c.readOpcode(c.pc))
	c.pc++
	return 0
}

// Zero-page addressing with X register offset
func (c *cpu6502) zpx() uint8 {
	c.addrAbs = uint16(c.readOpcode(c.pc) + c.x)
	c.pc++
	return 0
}

// Zero-page addressing with Y register offset
func (c *cpu6502) zpy() uint8 {
	c.addrAbs = uint16(c.readOpcode(c.pc) + c.y)
	c.pc++
	return 0
}

// Absolute addressing - full address in it's natural form
func (c *cpu6502) abs() uint8 {

}