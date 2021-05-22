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

	// Configure addressing modes
	abs := addressingMode{lblAddressingModeABS, c.abs}
	abx := addressingMode{lblAddressingModeABX, c.abx}
	aby := addressingMode{lblAddressingModeABY, c.aby}
	imp := addressingMode{lblAddressingModeIMP, c.imp}
	imm := addressingMode{lblAddressingModeIMM, c.imm}
	ind := addressingMode{lblAddressingModeIND, c.ind}
	izx := addressingMode{lblAddressingModeIZX, c.izx}
	izy := addressingMode{lblAddressingModeIZY, c.izy}
	zp0 := addressingMode{lblAddressingModeZP0, c.zp0}
	zpx := addressingMode{lblAddressingModeZPX, c.zpx}
	zpy := addressingMode{lblAddressingModeZPY, c.zpy}
	rel := addressingMode{lblAddressingModeREL, c.rel}

	// Configure instructions
	c.lookupOpcodes = []instruction{
		{lblOpcodeBRK, c.brk, imm, 7},
		{lblOpcodeORA, c.ora, izx, 6},
		{lblOpcodeXXX, c.xxx, imp, 2},
		{lblOpcodeXXX, c.xxx, imp, 8},
		{lblOpcodeXXX, c.nop, imp, 3},
		{lblOpcodeORA, c.ora, zp0, 3},
		{lblOpcodeASL,c.asl,zp0, 5},
		{lblOpcodeXXX,c.xxx,imp, 5},
		{lblOpcodePHP,c.php,imp, 3},
		{lblOpcodeORA,c.ora,imm, 2},
		{lblOpcodeASL,c.asl,imp, 2},
		{lblOpcodeXXX,c.xxx,imp, 2},
		{lblOpcodeXXX,c.nop,imp, 4},
		{lblOpcodeORA,c.ora,abs, 4},
		{lblOpcodeASL,c.asl,abs, 6},
		{lblOpcodeXXX,c.xxx,imp, 6},
		{lblOpcodeBPL,c.bpl,rel, 2},
		{lblOpcodeORA,c.ora,izy, 5},
		{lblOpcodeXXX,c.xxx,imp, 2},
		{lblOpcodeXXX,c.xxx,imp, 8},
		{lblOpcodeXXX,c.nop,imp, 4},
		{lblOpcodeORA,c.ora,zpx, 4},
		{lblOpcodeASL,c.asl,zpx, 6},
		{lblOpcodeXXX,c.xxx,imp, 6},
		{lblOpcodeCLC,c.clc,imp, 2},
		{lblOpcodeORA,c.ora,aby, 4},
		{lblOpcodeXXX,c.nop,imp, 2},
		{lblOpcodeXXX,c.xxx,imp, 7},
		{lblOpcodeXXX,c.nop,imp, 4},
		{lblOpcodeORA,c.ora,abx, 4},
		{lblOpcodeASL,c.asl,abx, 7},
		{lblOpcodeXXX,c.xxx,imp, 7},
		{lblOpcodeJSR,c.jsr,abs, 6},
		{lblOpcodeAND,c.and,izx, 6},
		{lblOpcodeXXX,c.xxx,imp, 2},
		{lblOpcodeXXX,c.xxx,imp, 8},
		{lblOpcodeBIT,c.bit,zp0, 3},
		{lblOpcodeAND,c.and,zp0, 3},
		{lblOpcodeROL,c.rol,zp0, 5},
		{lblOpcodeXXX,c.xxx,imp, 5},
		{lblOpcodePLP,c.plp,imp, 4},
		{lblOpcodeAND,c.and,imm, 2},
		{lblOpcodeROL,c.rol,imp, 2},
		{lblOpcodeXXX,c.xxx,imp, 2},
		{lblOpcodeBIT,c.bit,abs, 4},
		{lblOpcodeAND,c.and,abs, 4},
		{lblOpcodeROL,c.rol,abs, 6},
		{lblOpcodeXXX,c.xxx,imp, 6},
		{lblOpcodeBMI,c.bmi,rel, 2},
		{lblOpcodeAND,c.and,izy, 5},
		{lblOpcodeXXX,c.xxx,imp, 2},
		{lblOpcodeXXX,c.xxx,imp, 8},
		{lblOpcodeXXX,c.nop,imp, 4},
		{lblOpcodeAND,c.and,zpx, 4},
		{lblOpcodeROL,c.rol,zpx, 6},
		{lblOpcodeXXX,c.xxx,imp, 6},
		{lblOpcodeSEC,c.sec,imp, 2},
		{lblOpcodeAND,c.and,aby, 4},
		{lblOpcodeXXX,c.nop,imp, 2},
		{lblOpcodeXXX,c.xxx,imp, 7},
		{lblOpcodeXXX,c.nop,imp, 4},
		{lblOpcodeAND,c.and,abx, 4},
		{lblOpcodeROL,c.rol,abx, 7},
		{lblOpcodeXXX,c.xxx,imp, 7},
		{lblOpcodeRTI,c.rti,imp, 6},
		{lblOpcodeEOR,c.eor,izx, 6},
		{lblOpcodeXXX,c.xxx,imp, 2},
		{lblOpcodeXXX,c.xxx,imp, 8},
		{lblOpcodeXXX,c.nop,imp, 3},
		{lblOpcodeEOR,c.eor,zp0, 3},
		{lblOpcodeLSR,c.lsr,zp0, 5},
		{lblOpcodeXXX,c.xxx,imp, 5},
		{lblOpcodePHA,c.pha,imp, 3},
		{lblOpcodeEOR,c.eor,imm, 2},
		{lblOpcodeLSR,c.lsr,imp, 2},
		{lblOpcodeXXX,c.xxx,imp, 2},
		{lblOpcodeJMP,c.jmp,abs, 3},
		{lblOpcodeEOR,c.eor,abs, 4},
		{lblOpcodeLSR,c.lsr,abs, 6},
		{lblOpcodeXXX,c.xxx,imp, 6},
		{lblOpcodeBVC,c.bvc,rel, 2},
		{lblOpcodeEOR,c.eor,izy, 5},
		{lblOpcodeXXX,c.xxx,imp, 2},
		{lblOpcodeXXX,c.xxx,imp, 8},
		{lblOpcodeXXX,c.nop,imp, 4},
		{lblOpcodeEOR,c.eor,zpx, 4},
		{lblOpcodeLSR,c.lsr,zpx, 6},
		{lblOpcodeXXX,c.xxx,imp, 6},
		{lblOpcodeCLI,c.cli,imp, 2},
		{lblOpcodeEOR,c.eor,aby, 4},
		{lblOpcodeXXX,c.nop,imp, 2},
		{lblOpcodeXXX,c.xxx,imp, 7},
		{lblOpcodeXXX,c.nop,imp, 4},
		{lblOpcodeEOR,c.eor,abx, 4},
		{lblOpcodeLSR,c.lsr,abx, 7},
		{lblOpcodeXXX,c.xxx,imp, 7},
		{lblOpcodeRTS,c.rts,imp, 6},
		{lblOpcodeADC,c.adc,izx, 6},
		{lblOpcodeXXX,c.xxx,imp, 2},
		{lblOpcodeXXX,c.xxx,imp, 8},
		{lblOpcodeXXX,c.nop,imp, 3},
		{lblOpcodeADC,c.adc,zp0, 3},
		{lblOpcodeROR,c.ror,zp0, 5},
		{lblOpcodeXXX,c.xxx,imp, 5},
		{lblOpcodePLA,c.pla,imp, 4},
		{lblOpcodeADC,c.adc,imm, 2},
		{lblOpcodeROR,c.ror,imp, 2},
		{lblOpcodeXXX,c.xxx,imp, 2},
		{lblOpcodeJMP,c.jmp,ind, 5},
		{lblOpcodeADC,c.adc,abs, 4},
		{lblOpcodeROR,c.ror,abs, 6},
		{lblOpcodeXXX,c.xxx,imp, 6},
		{lblOpcodeBVS,c.bvs,rel, 2},
		{lblOpcodeADC,c.adc,izy, 5},
		{lblOpcodeXXX,c.xxx,imp, 2},
		{lblOpcodeXXX,c.xxx,imp, 8},
		{lblOpcodeXXX,c.nop,imp, 4},
		{lblOpcodeADC,c.adc,zpx, 4},
		{lblOpcodeROR,c.ror,zpx, 6},
		{lblOpcodeXXX,c.xxx,imp, 6},
		{lblOpcodeSEI,c.sei,imp, 2},
		{lblOpcodeADC,c.adc,aby, 4},
		{lblOpcodeXXX,c.nop,imp, 2},
		{lblOpcodeXXX,c.xxx,imp, 7},
		{lblOpcodeXXX,c.nop,imp, 4},
		{lblOpcodeADC,c.adc,abx, 4},
		{lblOpcodeROR,c.ror,abx, 7},
		{lblOpcodeXXX,c.xxx,imp, 7},
		{lblOpcodeXXX,c.nop,imp, 2},
		{lblOpcodeSTA,c.sta,izx, 6},
		{lblOpcodeXXX,c.nop,imp, 2},
		{lblOpcodeXXX,c.xxx,imp, 6},
		{lblOpcodeSTY,c.sty,zp0, 3},
		{lblOpcodeSTA,c.sta,zp0, 3},
		{lblOpcodeSTX,c.stx,zp0, 3},
		{lblOpcodeXXX,c.xxx,imp, 3},
		{lblOpcodeDEY,c.dey,imp, 2},
		{lblOpcodeXXX,c.nop,imp, 2},
		{lblOpcodeTXA,c.txa,imp, 2},
		{lblOpcodeXXX,c.xxx,imp, 2},
		{lblOpcodeSTY,c.sty,abs, 4},
		{lblOpcodeSTA,c.sta,abs, 4},
		{lblOpcodeSTX,c.stx,abs, 4},
		{lblOpcodeXXX,c.xxx,imp, 4},
		{lblOpcodeBCC,c.bcc,rel, 2},
		{lblOpcodeSTA,c.sta,izy, 6},
		{lblOpcodeXXX,c.xxx,imp, 2},
		{lblOpcodeXXX,c.xxx,imp, 6},
		{lblOpcodeSTY,c.sty,zpx, 4},
		{lblOpcodeSTA,c.sta,zpx, 4},
		{lblOpcodeSTX,c.stx,zpy, 4},
		{lblOpcodeXXX,c.xxx,imp, 4},
		{lblOpcodeTYA,c.tya,imp, 2},
		{lblOpcodeSTA,c.sta,aby, 5},
		{lblOpcodeTXS,c.txs,imp, 2},
		{lblOpcodeXXX,c.xxx,imp, 5},
		{lblOpcodeXXX,c.nop,imp, 5},
		{lblOpcodeSTA,c.sta,abx, 5},
		{lblOpcodeXXX,c.xxx,imp, 5},
		{lblOpcodeXXX,c.xxx,imp, 5},
		{lblOpcodeLDY,c.ldy,imm, 2},
		{lblOpcodeLDA,c.lda,izx, 6},
		{lblOpcodeLDX,c.ldx,imm, 2},
		{lblOpcodeXXX,c.xxx,imp, 6},
		{lblOpcodeLDY,c.ldy,zp0, 3},
		{lblOpcodeLDA,c.lda,zp0, 3},
		{lblOpcodeLDX,c.ldx,zp0, 3},
		{lblOpcodeXXX,c.xxx,imp, 3},
		{lblOpcodeTAY,c.tay,imp, 2},
		{lblOpcodeLDA,c.lda,imm, 2},
		{lblOpcodeTAX,c.tax,imp, 2},
		{lblOpcodeXXX,c.xxx,imp, 2},
		{lblOpcodeLDY,c.ldy,abs, 4},
		{lblOpcodeLDA,c.lda,abs, 4},
		{lblOpcodeLDX,c.ldx,abs, 4},
		{lblOpcodeXXX,c.xxx,imp, 4},
		{lblOpcodeBCS,c.bcs,rel, 2},
		{lblOpcodeLDA,c.lda,izy, 5},
		{lblOpcodeXXX,c.xxx,imp, 2},
		{lblOpcodeXXX,c.xxx,imp, 5},
		{lblOpcodeLDY,c.ldy,zpx, 4},
		{lblOpcodeLDA,c.lda,zpx, 4},
		{lblOpcodeLDX,c.ldx,zpy, 4},
		{lblOpcodeXXX,c.xxx,imp, 4},
		{lblOpcodeCLV,c.clv,imp, 2},
		{lblOpcodeLDA,c.lda,aby, 4},
		{lblOpcodeTSX,c.tsx,imp, 2},
		{lblOpcodeXXX,c.xxx,imp, 4},
		{lblOpcodeLDY,c.ldy,abx, 4},
		{lblOpcodeLDA,c.lda,abx, 4},
		{lblOpcodeLDX,c.ldx,aby, 4},
		{lblOpcodeXXX,c.xxx,imp, 4},
		{lblOpcodeCPY,c.cpy,imm, 2},
		{lblOpcodeCMP,c.cmp,izx, 6},
		{lblOpcodeXXX,c.nop,imp, 2},
		{lblOpcodeXXX,c.xxx,imp, 8},
		{lblOpcodeCPY,c.cpy,zp0, 3},
		{lblOpcodeCMP,c.cmp,zp0, 3},
		{lblOpcodeDEC,c.dec,zp0, 5},
		{lblOpcodeXXX,c.xxx,imp, 5},
		{lblOpcodeINY,c.iny,imp, 2},
		{lblOpcodeCMP,c.cmp,imm, 2},
		{lblOpcodeDEX,c.dex,imp, 2},
		{lblOpcodeXXX,c.xxx,imp, 2},
		{lblOpcodeCPY,c.cpy,abs, 4},
		{lblOpcodeCMP,c.cmp,abs, 4},
		{lblOpcodeDEC,c.dec,abs, 6},
		{lblOpcodeXXX,c.xxx,imp, 6},
		{lblOpcodeBNE,c.bne,rel, 2},
		{lblOpcodeCMP,c.cmp,izy, 5},
		{lblOpcodeXXX,c.xxx,imp, 2},
		{lblOpcodeXXX,c.xxx,imp, 8},
		{lblOpcodeXXX,c.nop,imp, 4},
		{lblOpcodeCMP,c.cmp,zpx, 4},
		{lblOpcodeDEC,c.dec,zpx, 6},
		{lblOpcodeXXX,c.xxx,imp, 6},
		{lblOpcodeCLD,c.cld,imp, 2},
		{lblOpcodeCMP,c.cmp,aby, 4},
		{lblOpcodeNOP,c.nop,imp, 2},
		{lblOpcodeXXX,c.xxx,imp, 7},
		{lblOpcodeXXX,c.nop,imp, 4},
		{lblOpcodeCMP,c.cmp,abx, 4},
		{lblOpcodeDEC,c.dec,abx, 7},
		{lblOpcodeXXX,c.xxx,imp, 7},
		{lblOpcodeCPX,c.cpx,imm, 2},
		{lblOpcodeSBC,c.sbc,izx, 6},
		{lblOpcodeXXX,c.nop,imp, 2},
		{lblOpcodeXXX,c.xxx,imp, 8},
		{lblOpcodeCPX,c.cpx,zp0, 3},
		{lblOpcodeSBC,c.sbc,zp0, 3},
		{lblOpcodeINC,c.inc,zp0, 5},
		{lblOpcodeXXX,c.xxx,imp, 5},
		{lblOpcodeINX,c.inx,imp, 2},
		{lblOpcodeSBC,c.sbc,imm, 2},
		{lblOpcodeNOP,c.nop,imp, 2},
		{lblOpcodeXXX,c.sbc,imp, 2},
		{lblOpcodeCPX,c.cpx,abs, 4},
		{lblOpcodeSBC,c.sbc,abs, 4},
		{lblOpcodeINC,c.inc,abs, 6},
		{lblOpcodeXXX,c.xxx,imp, 6},
		{lblOpcodeBEQ,c.beq,rel, 2},
		{lblOpcodeSBC,c.sbc,izy, 5},
		{lblOpcodeXXX,c.xxx,imp, 2},
		{lblOpcodeXXX,c.xxx,imp, 8},
		{lblOpcodeXXX,c.nop,imp, 4},
		{lblOpcodeSBC,c.sbc,zpx, 4},
		{lblOpcodeINC,c.inc,zpx, 6},
		{lblOpcodeXXX,c.xxx,imp, 6},
		{lblOpcodeSED,c.sed,imp, 2},
		{lblOpcodeSBC,c.sbc,aby, 4},
		{lblOpcodeNOP,c.nop,imp, 2},
		{lblOpcodeXXX,c.xxx,imp, 7},
		{lblOpcodeXXX,c.nop,imp, 4},
		{lblOpcodeSBC,c.sbc,abx, 4},
		{lblOpcodeINC,c.inc,abx, 7},
		{lblOpcodeXXX,c.xxx,imp, 7},
	}
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

func (c *cpu6502) fetchOpcode() uint8 {
	if c.lookupOpcodes[c.opcode].addressing.name == lblAddressingModeIMP {
		c.fetched = c.read(c.addrAbs)
	}
	return c.fetched
}
