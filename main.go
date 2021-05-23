package main

import (
	"github.com/mtojek/nes-emulator/bus"
	"github.com/mtojek/nes-emulator/cpu"
	"github.com/mtojek/nes-emulator/memory"
)

func main() {
	var b bus.Bus

	r := memory.Create()
	b.Connect(0x0000, 0x1FFF, r)

	cpu.Create(&b)
}
