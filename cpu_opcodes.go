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

// Break
func (c *cpu6502) brk() uint8 {
	c.pc++

	c.setFlag(flagI, true)
	c.bus.write(0x0100 + uint16(c.sp), uint8(c.pc >> 8))
	c.sp--
	c.bus.write(0x0100 + uint16(c.sp), uint8(c.pc))
	c.sp--

	c.setFlag(flagB, true)
	c.bus.write(0x0100 + uint16(c.sp), c.status)
	c.sp--
	c.setFlag(flagB, false)

	c.pc = uint16(c.bus.read(0xFFFE, true)) | uint16(c.bus.read(0xFFFF, true)) << 8
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

// Unknown instruction
func (c *cpu6502) xxx() uint8 {
	return 0
}
