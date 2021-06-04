package cpu

const (
	lblAddressingModeABS = "ABS"
	lblAddressingModeABX = "ABX"
	lblAddressingModeABY = "ABY"
	lblAddressingModeIMM = "IMM"
	lblAddressingModeIMP = "IMP"
	lblAddressingModeIND = "IND"
	lblAddressingModeIZX = "IZX"
	lblAddressingModeIZY = "IZY"
	lblAddressingModeREL = "REL"
	lblAddressingModeZP0 = "ZP0"
	lblAddressingModeZPX = "ZPX"
	lblAddressingModeZPY = "ZPY"
)

type addressingMode struct {
	name string
	do   addressingModeFunc
}

type addressingModeFunc func() uint8

// Implied addressing mode
func (c *CPU6502) imp() uint8 {
	c.fetched = c.a
	return 0
}

// Immediate addressing mode
func (c *CPU6502) imm() uint8 {
	c.addrAbs = c.pc
	c.pc++
	return 0
}

// Zero-page addressing
func (c *CPU6502) zp0() uint8 {
	c.addrAbs = uint16(c.read(c.pc))
	c.pc++
	return 0
}

// Zero-page addressing with X register offset
func (c *CPU6502) zpx() uint8 {
	c.addrAbs = uint16(c.read(c.pc) + c.x)
	c.pc++
	return 0
}

// Zero-page addressing with Y register offset
func (c *CPU6502) zpy() uint8 {
	c.addrAbs = uint16(c.read(c.pc) + c.y)
	c.pc++
	return 0
}

// Absolute addressing with X register offset
func (c *CPU6502) abx() uint8 {
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
func (c *CPU6502) aby() uint8 {
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
func (c *CPU6502) ind() uint8 {
	ptrLo := uint16(c.read(c.pc))
	c.pc++

	ptrHi := uint16(c.read(c.pc)) << 8
	c.pc++

	ptr := ptrHi | ptrLo

	if ptrLo == 0x00FF { // Simulate page boundary hardware bug
		c.addrAbs = (uint16(c.read(ptr&0xFF00)) << 8) | uint16(c.read(ptr))
		return 0
	}

	// Behave normally
	c.addrAbs = (uint16(c.read(ptr+1)) << 8) | uint16(c.read(ptr))
	return 0
}

// Indirect addressing of the zero-page with X offset
func (c *CPU6502) izx() uint8 {
	t := uint16(c.read(c.pc))
	c.pc++

	lo := uint16(c.read((t + uint16(c.x)) & 0x00FF))
	hi := uint16(c.read((t+uint16(c.x)+1)&0x00FF)) << 8
	c.addrAbs = hi | lo
	return 0
}

// Indirect addressing of the zero-page with Y offset
func (c *CPU6502) izy() uint8 {
	address := uint16(c.read(c.pc))
	c.pc++

	a := address
	b := (a & 0xFF00) | uint16(byte(a)+1)
	lo := c.read(a)
	hi := c.read(b)
	c.addrAbs = (uint16(hi)<<8 | uint16(lo)) + uint16(c.y)
	if (c.addrAbs-uint16(c.y))&0xFF00 != c.addrAbs&0xFF00 {
		return 1
	}
	return 0
}

// Relative addressing mode
func (c *CPU6502) rel() uint8 {
	c.addrRel = uint16(c.read(c.pc))
	c.pc++
	if c.addrRel&0x80 == 0x80 {
		c.addrRel |= 0xFF00
	}
	return 0
}

// Absolute addressing mode
func (c *CPU6502) abs() uint8 {
	lo := uint16(c.read(c.pc))
	c.pc++

	hi := uint16(c.read(c.pc)) << 8
	c.pc++

	c.addrAbs = hi + lo
	return 0
}
