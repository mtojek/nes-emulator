package memory

import "github.com/mtojek/nes-emulator/bus"

type Mirroring struct {
	destination bus.ReadableWriteable

	startOffset uint16
	mask        uint16
}

var _ bus.ReadableWriteable = new(Mirroring)

func CreateMirroring(destination bus.ReadableWriteable, startOffset, mask uint16) *Mirroring {
	return &Mirroring{
		destination: destination,
		startOffset: startOffset,
		mask:        mask,
	}
}

func (m Mirroring) Read(addr uint16) uint8 {
	return m.destination.Read(m.startOffset + (addr & m.mask))
}

func (m Mirroring) Write(addr uint16, data uint8) {
	m.destination.Write(m.startOffset+(addr&m.mask), data)
}
