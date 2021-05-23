package ram

import "github.com/mtojek/nes-emulator/bus"

const defaultRAMSize = 2048

type RAM struct {
	memory []uint8
}

var _ bus.ReadableWriteable = new(RAM)

func Create() *RAM {
	return CreateWithSize(defaultRAMSize)
}

func CreateWithSize(size uint16) *RAM {
	return &RAM{
		memory: make([]uint8, size),
	}
}

func (r *RAM) Write(addr uint16, data uint8) {
	r.memory[addr] = data
}

func (r *RAM) Read(addr uint16, _ bool) uint8 {
	return r.memory[addr]
}
