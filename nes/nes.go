package nes

import (
	"github.com/mtojek/nes-emulator/bus"
	"github.com/mtojek/nes-emulator/cpu"
	"github.com/mtojek/nes-emulator/memory"
)

type NES struct{}

func Create() *NES {
	var b bus.Bus

	ram := memory.CreateMemory()
	mirroredRAM := memory.CreateMirroring(ram, 0x07FF)
	b.Connect(0x0000, 0x1FFF, mirroredRAM)

	cpu.Create(&b)

	return new(NES)
}
