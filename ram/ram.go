package ram

import "github.com/mtojek/nes-emulator/bus"

type RAM struct {
	memory []uint8
}

var _ bus.ReadableWriteable = new(RAM)

func Create() *RAM {
	return &RAM{
		memory: make([]uint8, 64*1024),
	}
}

func (r *RAM) Write(addr uint16, data uint8) {
	r.memory[addr] = data
}

func (r *RAM) Read(addr uint16, _ bool) uint8 {
	return r.memory[addr]
}
