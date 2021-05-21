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
	c.addrAbs = uint16(c.read(c.pc))
	c.pc++
	return 0
}

// Zero-page addressing with X register offset
func (c *cpu6502) zpx() uint8 {
	c.addrAbs = uint16(c.read(c.pc) + c.x)
	c.pc++
	return 0
}

// Zero-page addressing with Y register offset
func (c *cpu6502) zpy() uint8 {
	c.addrAbs = uint16(c.read(c.pc) + c.y)
	c.pc++
	return 0
}

// Absolute addressing with X register offset
func (c *cpu6502) abx() uint8 {
	lo := uint16(c.read(c.pc))
	c.pc++

	hi := uint16(c.read(c.pc)) << 8
	c.pc++

	c.addrAbs = hi + lo + uint16(c.x)

	if c.addrAbs&0xFF00 != hi {
		return 1
	}
	return 0
}

// Absolute addressing with Y register offset
func (c *cpu6502) aby() uint8 {
	lo := uint16(c.read(c.pc))
	c.pc++

	hi := uint16(c.read(c.pc)) << 8
	c.pc++

	c.addrAbs = hi + lo + uint16(c.y)

	if c.addrAbs&0xFF00 != hi {
		return 1
	}
	return 0
}

// Indirect addressing
func (c *cpu6502) ind() uint8 {
	ptrLo := uint16(c.read(c.pc))
	c.pc++

	ptrHi := uint16(c.read(c.pc)) << 8
	c.pc++

	ptr := ptrHi + ptrLo

	if ptrLo == 0x00FF { // Simulate page boundary hardware bug
		c.addrAbs = (uint16(c.read(ptr&0x00FF)) << 8) | uint16(c.read(ptr))
		return 0
	}

	// Behave normally
	c.addrAbs = (uint16(c.read(ptr+1)) << 8) | uint16(c.read(ptr))
	return 0
}

// Indirect addressing of the zero-page with X offset
func (c *cpu6502) izx() uint8 {
	t := uint16(c.read(c.pc))
	c.pc++

	lo := uint16(c.read((t + uint16(c.x)) & 0x00FF))
	hi := uint16(c.read((t+uint16(c.x)+1)&0x00FF)) << 8
	c.addrAbs = hi + lo
	return 0
}

// Indirect addressing of the zero-page with Y offset
func (c *cpu6502) izy() uint8 {
	t := uint16(c.read(c.pc))
	c.pc++

	lo := uint16(c.read(t & 0x00FF))
	hi := uint16(c.read((t+1)&0x00FF)) << 8
	c.addrAbs = hi + lo + uint16(c.y)

	if c.addrAbs&0xFF00 != hi {
		return 1
	}
	return 0
}
