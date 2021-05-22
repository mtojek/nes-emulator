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
	name   string

	operationDo operateFunc
	addressing  addressingMode

	cycles uint8
}

type operateFunc func() uint8

// Bitwise logic AND
func (c *cpu6502) and() uint8 {
	c.fetchOpcode()
	c.a = c.a & c.fetched
	c.setFlag(flagZ, c.a == 0x00)
	c.setFlag(flagN, (c.a&0x80) == 0x80)
	return 1
}

// Unknown instruction
func (c *cpu6502) xxx() uint8 {
	return 0
}
