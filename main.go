package main

import (
	"github.com/mtojek/nes-emulator/bus"
	"github.com/mtojek/nes-emulator/cpu"
	"github.com/mtojek/nes-emulator/ram"
)

func main() {
	var b bus.Bus

	r := ram.Create()
	b.Connect(0x0000, 0x1FFF, r)

	cpu.Create(&b)
}
