package main

const (
	lblOpcodeADC = "ADC"
	lblOpcodeAND = "AND"
	lblOpcodeASL = "ASL"
	lblOpcodeBCC = "BCC"
	lblOpcodeBCS = "BCS"
	lblOpcodeBEQ = "BEQ"
	lblOpcodeBIT = "BIT"
	lblOpcodeBMI = "BMI"
	lblOpcodeBNE = "BNE"
	lblOpcodeBPL = "BPL"
	lblOpcodeBRK = "BRK"
	lblOpcodeBVC = "BVC"
	lblOpcodeBVS = "BVS"
	lblOpcodeCLC = "CLC"
	lblOpcodeCLD = "CLD"
	lblOpcodeCLI = "CLI"
	lblOpcodeCLV = "CLV"
	lblOpcodeCMP = "CMP"
	lblOpcodeCPX = "CPX"
	lblOpcodeCPY = "CPY"
	lblOpcodeDEC = "DEC"
	lblOpcodeDEX = "DEX"
	lblOpcodeDEY = "DEY"
	lblOpcodeEOR = "EOR"
	lblOpcodeINC = "INC"
	lblOpcodeINX = "INX"
	lblOpcodeINY = "INY"
	lblOpcodeJMP = "JMP"
	lblOpcodeJSR = "JSR"
	lblOpcodeLDA = "LDA"
	lblOpcodeLDX = "LDX"
	lblOpcodeLDY = "LDY"
	lblOpcodeLSR = "LSR"
	lblOpcodeNOP = "NOP"
	lblOpcodeORA = "ORA"
	lblOpcodePHA = "PHA"
	lblOpcodePHP = "PHP"
	lblOpcodePLA = "PLA"
	lblOpcodePLP = "PLP"
	lblOpcodeROL = "ROL"
	lblOpcodeROR = "ROR"
	lblOpcodeRTI = "RTI"
	lblOpcodeRTS = "RTS"
	lblOpcodeSBC = "SBC"
	lblOpcodeSEC = "SEC"
	lblOpcodeSED = "SED"
	lblOpcodeSEI = "SEI"
	lblOpcodeSTA = "STA"
	lblOpcodeSTX = "STX"
	lblOpcodeSTY = "STY"
	lblOpcodeTAY = "TAY"
	lblOpcodeTAX = "TAX"
	lblOpcodeTSX = "TSX"
	lblOpcodeTXA = "TXA"
	lblOpcodeTXS = "TXS"
	lblOpcodeTYA = "TYA"
	lblOpcodeXXX = "???"
)

type instruction struct {
	name string

	operationDo operateFunc
	addressing  addressingMode

	cycles uint8
}

type operateFunc func() uint8

// Add with Carry In
func (c *cpu6502) adc() uint8 {
	c.fetch()
	t := uint16(c.a) + uint16(c.fetched) + uint16(c.getFlag(flagC))
	c.setFlag(flagC, t > 255)
	c.setFlag(flagV, (uint16(c.a)^uint16(c.fetched))&0x0080 == 0 && (uint16(c.a)^t)&0x0080 != 0)
	c.setFlagZ(uint8(t))
	c.setFlagN(uint8(t))
	c.a = uint8(t)
	return 1
}

// Bitwise logic AND
func (c *cpu6502) and() uint8 {
	c.fetch()
	c.a = c.a & c.fetched
	c.setFlagZ(c.a)
	c.setFlagN(c.a)
	return 1
}

// Arithmetic Shift Left
func (c *cpu6502) asl() uint8 {
	c.fetch()
	t := uint16(c.fetched) << 1
	c.setFlag(flagC, (t&0xFF00) > 0)
	c.setFlagZ(uint8(t))
	c.setFlagN(uint8(t))
	if c.lookupOpcodes[c.opcode].addressing.name == lblAddressingModeIMP {
		c.a = uint8(t)
	} else {
		c.bus.write(c.addrAbs, uint8(t))
	}
	return 0
}

// Branch if Carry Clear
func (c *cpu6502) bcc() uint8 {
	if c.getFlag(flagC) == 0 {
		c.cycles++
		c.addrAbs = c.pc + c.addrRel

		if (c.addrAbs & 0xFF00) != (c.pc & 0xFF00) {
			c.cycles++
		}
		c.pc = c.addrAbs
	}
	return 0
}

// Branch if Carry Set
func (c *cpu6502) bcs() uint8 {
	if c.getFlag(flagC) == 1 {
		c.cycles++
		c.addrAbs = c.pc + c.addrRel

		if (c.addrAbs & 0xFF00) != (c.pc & 0xFF00) {
			c.cycles++
		}
		c.pc = c.addrAbs
	}
	return 0
}

// Branch if Equal
func (c *cpu6502) beq() uint8 {
	if c.getFlag(flagZ) == 1 {
		c.cycles++
		c.addrAbs = c.pc + c.addrRel

		if (c.addrAbs & 0xFF00) != (c.pc & 0xFF00) {
			c.cycles++
		}
		c.pc = c.addrAbs
	}
	return 0
}

// Branch if Not Equal
func (c *cpu6502) bne() uint8 {
	if c.getFlag(flagZ) == 0 {
		c.cycles++
		c.addrAbs = c.pc + c.addrRel

		if (c.addrAbs & 0xFF00) != (c.pc & 0xFF00) {
			c.cycles++
		}
		c.pc = c.addrAbs
	}
	return 0
}

// Branch if Overflow Set
func (c *cpu6502) bvs() uint8 {
	if c.getFlag(flagV) == 1 {
		c.cycles++
		c.addrAbs = c.pc + c.addrRel

		if (c.addrAbs & 0xFF00) != (c.pc & 0xFF00) {
			c.cycles++
		}
		c.pc = c.addrAbs
	}
	return 0
}

// Test bit
func (c *cpu6502) bit() uint8 {
	c.fetch()
	t := c.a & c.fetched
	c.setFlagZ(t)
	c.setFlagN(c.fetched)
	c.setFlag(flagV, (c.fetched&flagV) == flagV)
	return 0
}

// Branch if Negative
func (c *cpu6502) bmi() uint8 {
	if c.getFlag(flagN) == 1 {
		c.cycles++
		c.addrAbs = c.pc + c.addrRel

		if (c.addrAbs & 0xFF00) != (c.pc & 0xFF00) {
			c.cycles++
		}
		c.pc = c.addrAbs
	}
	return 0
}

// Branch if Positive
func (c *cpu6502) bpl() uint8 {
	if c.getFlag(flagN) == 0 {
		c.cycles++
		c.addrAbs = c.pc + c.addrRel

		if (c.addrAbs & 0xFF00) != (c.pc & 0xFF00) {
			c.cycles++
		}
		c.pc = c.addrAbs
	}
	return 0
}

// Break
func (c *cpu6502) brk() uint8 {
	c.pc++

	c.setFlag(flagI, true)
	c.bus.write(0x0100+uint16(c.sp), uint8(c.pc>>8))
	c.sp--
	c.bus.write(0x0100+uint16(c.sp), uint8(c.pc))
	c.sp--

	c.setFlag(flagB, true)
	c.bus.write(0x0100+uint16(c.sp), c.status)
	c.sp--
	c.setFlag(flagB, false)

	c.pc = uint16(c.bus.read(0xFFFE, true)) | uint16(c.bus.read(0xFFFF, true))<<8
	return 0
}

// Branch if Overflow Clear
func (c *cpu6502) bvc() uint8 {
	if c.getFlag(flagV) == 0 {
		c.cycles++
		c.addrAbs = c.pc + c.addrRel

		if (c.addrAbs & 0xFF00) != (c.pc & 0xFF00) {
			c.cycles++
		}
		c.pc = c.addrAbs
	}
	return 0
}

// Clear Carry Flag
func (c *cpu6502) clc() uint8 {
	c.setFlag(flagC, false)
	return 0
}

// Clear Decimal Flag
func (c *cpu6502) cld() uint8 {
	c.setFlag(flagD, false)
	return 0
}

// Disable Interrupts / Clear Interrupt Flag
func (c *cpu6502) cli() uint8 {
	c.setFlag(flagI, false)
	return 0
}

// Clear Overflow Flag
func (c *cpu6502) clv() uint8 {
	c.setFlag(flagV, false)
	return 0
}

// Compare Accumulator
func (c *cpu6502) cmp() uint8 {
	c.fetch()
	t := uint16(c.a) - uint16(c.fetched)
	c.setFlag(flagC, c.a >= c.fetched)
	c.setFlagZ(uint8(t))
	c.setFlagN(uint8(t))
	return 1
}

// Compare X Register
func (c *cpu6502) cpx() uint8 {
	c.fetch()
	t := uint16(c.x) - uint16(c.fetched)
	c.setFlag(flagC, c.x >= c.fetched)
	c.setFlagZ(uint8(t))
	c.setFlagN(uint8(t))
	return 0
}

// Compare Y Register
func (c *cpu6502) cpy() uint8 {
	c.fetch()
	t := uint16(c.y) - uint16(c.fetched)
	c.setFlag(flagC, c.y >= c.fetched)
	c.setFlagZ(uint8(t))
	c.setFlagN(uint8(t))
	return 0
}

// Decrement Value at Memory Location
func (c *cpu6502) dec() uint8 {
	c.fetch()
	t := c.fetched - 1
	c.bus.write(c.addrAbs, t&0x00FF)
	c.setFlagZ(t)
	c.setFlagN(t)
	return 0
}

// Decrement X Register
func (c *cpu6502) dex() uint8 {
	c.x--
	c.setFlagZ(c.x)
	c.setFlagN(c.x)
	return 0
}

// Decrement Y Register
func (c *cpu6502) dey() uint8 {
	c.y--
	c.setFlagZ(c.y)
	c.setFlagN(c.y)
	return 0
}

// Bitwise Logic XOR
func (c *cpu6502) eor() uint8 {
	c.fetch()
	c.a = c.a ^ c.fetched
	c.setFlagZ(c.a)
	c.setFlagN(c.a)
	return 1
}

// Increment Value at Memory Location
func (c *cpu6502) inc() uint8 {
	c.fetch()
	t := uint16(c.fetched) + 1
	c.bus.write(c.addrAbs, uint8(t))
	c.setFlagZ(uint8(t))
	c.setFlagN(uint8(t))
	return 0
}

// Increment X Register
func (c *cpu6502) inx() uint8 {
	c.x++
	c.setFlagZ(c.x)
	c.setFlagN(c.x)
	return 0
}

// Increment Y Register
func (c *cpu6502) iny() uint8 {
	c.y++
	c.setFlagZ(c.y)
	c.setFlagN(c.y)
	return 0
}

// Jump To Location
func (c *cpu6502) jmp() uint8 {
	c.pc = c.addrAbs
	return 0
}

// Jump To Sub-Routine
func (c *cpu6502) jsr() uint8 {
	c.pc--

	c.bus.write(0x0100+uint16(c.sp), uint8(c.pc>>8))
	c.sp--
	c.bus.write(0x0100+uint16(c.sp), uint8(c.pc))
	c.sp--

	c.pc = c.addrAbs
	return 0
}

// Load The Accumulator
func (c *cpu6502) lda() uint8 {
	c.fetch()
	c.a = c.fetched
	c.setFlagZ(c.a)
	c.setFlagN(c.a)
	return 1
}

// Load The X Register
func (c *cpu6502) ldx() uint8 {
	c.fetch()
	c.x = c.fetched
	c.setFlagZ(c.x)
	c.setFlagN(c.x)
	return 1
}

// Load The Y Register
func (c *cpu6502) ldy() uint8 {
	c.fetch()
	c.y = c.fetched
	c.setFlagZ(c.y)
	c.setFlagN(c.y)
	return 1
}

// Logical Shift Right
func (c *cpu6502) lsr() uint8 {
	c.fetch()
	c.setFlag(flagC, c.fetched&flagC == 1)
	t := c.fetched >> 1
	c.setFlagZ(t)
	c.setFlagN(t)
	if c.lookupOpcodes[c.opcode].addressing.name == lblAddressingModeIMP {
		c.a = t
	} else {
		c.bus.write(c.addrAbs, t)
	}
	return 0
}

// No Operation
func (c *cpu6502) nop() uint8 {
	// Illegal opcodes
	// See reference: https://wiki.nesdev.com/w/index.php/CPU_unofficial_opcodes
	switch c.opcode {
	case 0x1C:
	case 0x3C:
	case 0x5C:
	case 0x7C:
	case 0xDC:
	case 0xFC:
		return 1
	}
	return 0
}

// Bitwise Logic OR
func (c *cpu6502) ora() uint8 {
	c.fetch()
	c.a = c.a | c.fetched
	c.setFlagZ(c.a)
	c.setFlagN(c.a)
	return 1
}

// Push Accumulator to Stack
func (c *cpu6502) pha() uint8 {
	c.bus.write(0x0100+uint16(c.sp), c.a)
	c.sp--
	return 0
}

// Pop Accumulator from Stack
func (c *cpu6502) pla() uint8 {
	c.sp++
	c.a = c.bus.read(0x0100+uint16(c.sp), true)
	c.setFlagZ(c.a)
	c.setFlagN(c.a)
	return 0
}

// Push Status Register to Stack
func (c *cpu6502) php() uint8 {
	c.bus.write(0x0100+uint16(c.sp), c.status|flagB|flagU)
	c.setFlag(flagB, false)
	c.setFlag(flagU, false)
	c.sp--
	return 0
}

// Pop Status Register off Stack
func (c *cpu6502) plp() uint8 {
	c.sp++
	c.status = c.bus.read(0x0100+uint16(c.sp), true)
	c.setFlag(flagU, true)
	return 0
}

// Rotate Left
func (c *cpu6502) rol() uint8 {
	c.fetch()
	t := uint16(c.fetched)<<1 | uint16(c.getFlag(flagC))
	c.setFlag(flagC, t&0xFF00 > 0)
	c.setFlagZ(uint8(t))
	c.setFlagN(uint8(t))
	if c.lookupOpcodes[c.opcode].addressing.name == lblAddressingModeIMP {
		c.a = uint8(t)
	} else {
		c.bus.write(c.addrAbs, uint8(t))
	}
	return 0
}

// Rotate Right
func (c *cpu6502) ror() uint8 {
	c.fetch()
	t := uint16(c.getFlag(flagC))<<7 | uint16(c.fetched)>>1
	c.setFlag(flagC, c.fetched&flagC == flagC)
	c.setFlagZ(uint8(t))
	c.setFlagN(uint8(t))
	if c.lookupOpcodes[c.opcode].addressing.name == "IMP" {
		c.a = uint8(t)
	} else {
		c.bus.write(c.addrAbs, uint8(t))
	}
	return 0
}

// Return from Interrupt
func (c *cpu6502) rti() uint8 {
	c.sp++
	c.status = c.bus.read(0x0100+uint16(c.sp), true)
	c.status &= ^flagB
	c.status &= ^flagU

	c.sp++
	c.pc = uint16(c.bus.read(0x0100+uint16(c.sp), true))
	c.sp++
	c.pc |= uint16(c.bus.read(0x0100+uint16(c.sp), true)) << 8
	return 0
}

// Return from Subroutine
func (c *cpu6502) rts() uint8 {
	c.sp++
	c.pc = uint16(c.bus.read(0x0100+uint16(c.sp), true))
	c.sp++
	c.pc |= uint16(c.bus.read(0x0100+uint16(c.sp), true)) << 8

	c.pc++
	return 0
}

// Subtract with Borrow In
func (c *cpu6502) sbc() uint8 {
	c.fetch()
	value := uint16(c.fetched) ^ 0x00FF
	t := uint16(c.a) + value + uint16(c.getFlag(flagC))
	c.setFlag(flagC, t > 255)
	c.setFlag(flagV, (uint16(c.a)^uint16(c.fetched))&0x0080 == 0 && (uint16(c.a)^t)&0x0080 != 0)
	c.setFlagZ(uint8(t))
	c.setFlagN(uint8(t))
	c.a = uint8(t)
	return 1
}

// Set Carry Flag
func (c *cpu6502) sec() uint8 {
	c.setFlag(flagC, true)
	return 0
}

// Set Carry Flag
func (c *cpu6502) sed() uint8 {
	c.setFlag(flagD, true)
	return 0
}

// Set Interrupt Flag
func (c *cpu6502) sei() uint8 {
	c.setFlag(flagI, true)
	return 0
}

// Store Accumulator at Address
func (c *cpu6502) sta() uint8 {
	c.bus.write(c.addrAbs, c.a)
	return 0
}

// Store X Register at Address
func (c *cpu6502) stx() uint8 {
	c.bus.write(c.addrAbs, c.x)
	return 0
}

// Store Y Register at Address
func (c *cpu6502) sty() uint8 {
	c.bus.write(c.addrAbs, c.y)
	return 0
}

// Transfer Stack Pointer to X Register
func (c *cpu6502) tsx() uint8 {
	c.x = c.sp
	c.setFlagZ(c.x)
	c.setFlagN(c.x)
	return 0
}

// Transfer Accumulator to X register
func (c *cpu6502) tax() uint8 {
	c.x = c.a
	c.setFlagZ(c.x)
	c.setFlagN(c.x)
	return 0
}

// Transfer Accumulator to Y register
func (c *cpu6502) tay() uint8 {
	c.y = c.a
	c.setFlagZ(c.y)
	c.setFlagN(c.y)
	return 0
}

// Transfer X register to Accumulator
func (c *cpu6502) txa() uint8 {
	c.a = c.x
	c.setFlagZ(c.a)
	c.setFlagN(c.a)
	return 0
}

// Transfer X register to Stack Pointer
func (c *cpu6502) txs() uint8 {
	c.sp = c.x
	return 0
}

// Transfer Y register to Accumulator
func (c *cpu6502) tya() uint8 {
	c.a = c.y
	c.setFlagZ(c.a)
	c.setFlagN(c.a)
	return 0
}

// Unknown instruction
func (c *cpu6502) xxx() uint8 {
	return 0
}
