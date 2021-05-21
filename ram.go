package main

type ram struct {
	memory []uint8
}

var _ readableWriteable = new(ram)

func createRAM() *ram {
	return &ram{
		memory: make([]uint8, 64*1024),
	}
}

func (r *ram) write(addr uint16, data uint8) {
	r.memory[addr] = data
}

func (r *ram) read(addr uint16, _ bool) uint8 {
	return r.memory[addr]
}
