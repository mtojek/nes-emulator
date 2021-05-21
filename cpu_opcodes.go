package main

type instruction struct {
	name   string
	cycles uint8

	addressingMode string
	addressingModeDo   addressingModeFunc

	operation string
	operationDo   operateFunc
}

type operateFunc func() uint8

func (c *cpu6502) fetch() uint8 {
	if c.lookupOpcodes[c.opcode].addressingMode == addressingModeIMP {
		c.fetched = c.read(c.addrAbs)
	}
	return c.fetched
}
