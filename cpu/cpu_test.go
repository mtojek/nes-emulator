package cpu_test

import (
	"encoding/hex"
	"github.com/mtojek/nes-emulator/bus"
	"github.com/mtojek/nes-emulator/cpu"
	"github.com/mtojek/nes-emulator/ram"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

const standardCodeLocation = 0x8000

/*
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
	var b bus.Bus

	r := ram.Create()
	b.Connect(0x0000, 0xFFFF, r)
	loadIntoRAM(t, r, standardCodeLocation, basicCode)
	setResetVector(r, standardCodeLocation)

	c := cpu.Create(&b)

	// when
	for i := 0; i < 1000; i++ {
		c.Clock()
	}
	b.Print(0x0000, 0x00FF)
	b.Print(standardCodeLocation, standardCodeLocation+0x00FF)

	// then
	require.Equal(t, r.Read(0x0000, true), uint8(0x0A))
	require.Equal(t, r.Read(0x0001, true), uint8(0x03))
	require.Equal(t, r.Read(0x0002, true), uint8(0x1E))
}

func loadIntoRAM(t *testing.T, ram bus.Writeable, offset uint16, code string) {
	code = strings.ReplaceAll(code, " ", "")
	code = strings.ReplaceAll(code, "\n", "")
	decoded, err := hex.DecodeString(code)
	require.NoError(t, err, "can't decode machine code")

	for i := uint16(0); i < uint16(len(decoded)); i++ {
		ram.Write(offset+i, decoded[i])
	}
}

func setResetVector(ram bus.Writeable, offset uint16) {
	ram.Write(0xFFFC, uint8(offset))
	ram.Write(0xFFFD, uint8(offset>>8))
}
