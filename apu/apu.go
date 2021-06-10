package apu

import (
	"fmt"

	"github.com/mtojek/nes-emulator/bus"
)

type APU2303 struct {
	pulse1 Pulse
	pulse2 Pulse
	noise Noise
	triangle Triangle
	dmc DMC
}

var _ bus.ReadableWriteable = new(APU2303)

func Create() *APU2303 {
	return new(APU2303)
}

func (apu *APU2303) Read(addr uint16) uint8 {
	switch addr {
	case 0x4015:
		return apu.readStatus()
	}
	panic(fmt.Sprintf("APU: read from unmapped address %04X", addr))
}

func (apu *APU2303) Write(addr uint16, data uint8) {
	panic("implement me")
}

func (apu *APU2303) readStatus() byte {
	var result byte
	if apu.pulse1.lengthValue > 0 {
		result |= 1
	}
	if apu.pulse2.lengthValue > 0 {
		result |= 2
	}
	if apu.triangle.lengthValue > 0 {
		result |= 4
	}
	if apu.noise.lengthValue > 0 {
		result |= 8
	}
	if apu.dmc.currentLength > 0 {
		result |= 16
	}
	return result
}