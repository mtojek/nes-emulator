package main

import (
	"encoding/hex"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

/**
*=$8000
LDX #10
STX $0000
LDX #3
STX $0001
LDY $0000
LDA #0
CLC
loop
ADC $0001
DEY
BNE loop
STA $002
NOP
NOP
NOP
*/
const basicCode = `
A2 0A 8E 00 00 A2 03 8E
01 00 AC 00 00 A9 00 18
6D 01 00 88 D0 FA 8D 02
00 EA EA EA
`

func TestCPU_BasicCode(t *testing.T) {
	// given
	var b bus

	r := createRAM()
	b.connect(0x0000, 0xFFFF, r)
	loadIntoRAM(t, r, 0x8000, basicCode)
	setResetVector(r, 0x8000)

	c := createCPU(&b)
	c.reset()

	// when
	for i := 0; i < 1000; i++ {
		c.clock()
	}

	// then
	require.Equal(t, r.memory[0x0000], uint8(0x0A))
	require.Equal(t, r.memory[0x0001], uint8(0x03))
	require.Equal(t, r.memory[0x0002], uint8(0x1E))
}

func loadIntoRAM(t *testing.T, ram writeable, offset uint16, code string) {
	code = strings.ReplaceAll(code, " ", "")
	code = strings.ReplaceAll(code, "\n", "")
	decoded, err := hex.DecodeString(code)
	require.NoError(t, err, "can't decode machine code")

	for i := uint16(0); i < uint16(len(decoded)); i++ {
		ram.write(offset+i, decoded[i])
	}
}

func setResetVector(ram writeable, offset uint16) {
	ram.write(0xFFFC, uint8(offset))
	ram.write(0xFFFD, uint8(offset>>8))
}
