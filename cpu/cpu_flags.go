package cpu

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

func (c *CPU6502) getFlag(f uint8) uint8 {
	if (c.status & f) > 0 {
		return 1
	}
	return 0
}

func (c *CPU6502) setFlagN(value uint8) {
	c.setFlag(flagN, (value&flagN) == flagN)
}

func (c *CPU6502) setFlagZ(value uint8) {
	c.setFlag(flagZ, value == 0x00)
}

func (c *CPU6502) setFlag(f uint8, value bool) {
	if value {
		c.status |= f
		return
	}
	c.status &= ^f
}
