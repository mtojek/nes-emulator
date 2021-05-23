package memory

import "github.com/mtojek/nes-emulator/bus"

const defaultRAMSize = 0x800

type Memory struct {
	blob []uint8
}

var _ bus.ReadableWriteable = new(Memory)

func Create() *Memory {
	return CreateWithSize(defaultRAMSize)
}

func CreateWithSize(size uint16) *Memory {
	return &Memory{
		blob: make([]uint8, size),
	}
}

func (r *Memory) Write(addr uint16, data uint8) {
	r.blob[addr] = data
}

func (r *Memory) Read(addr uint16, _ bool) uint8 {
	return r.blob[addr]
}
