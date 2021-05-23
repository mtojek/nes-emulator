package memory

import "github.com/mtojek/nes-emulator/bus"

type Mirroring struct {
	destination bus.ReadableWriteable
	mask        uint16
}

var _ bus.ReadableWriteable = new(Mirroring)

func CreateMirroring(destination bus.ReadableWriteable, mask uint16) *Mirroring {
	return &Mirroring{
		destination: destination,
		mask:        mask,
	}
}

func (m Mirroring) Read(addr uint16, bReadOnly bool) uint8 {
	return m.destination.Read(addr&m.mask, bReadOnly)
}

func (m Mirroring) Write(addr uint16, data uint8) {
	m.destination.Write(addr&m.mask, data)
}
