package memory

import "github.com/mtojek/nes-emulator/bus"

const (
	defaultRAMSize     = 0x0800
	defaultStartOffset = 0x0000
)

type Memory struct {
	startOffset uint16
	blob        []uint8
}

var _ bus.ReadableWriteable = new(Memory)

func CreateMemory() *Memory {
	return CreateMemoryWithSize(defaultStartOffset, defaultRAMSize)
}

func CreateMemoryWithSize(startOffset, size uint16) *Memory {
	return &Memory{
		startOffset: startOffset,
		blob:        make([]uint8, size),
	}
}

func (r *Memory) Write(addr uint16, data uint8) {
	r.blob[addr-r.startOffset] = data
}

func (r *Memory) Read(addr uint16) uint8 {
	return r.blob[addr-r.startOffset]
}
