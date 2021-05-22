package main

import (
	"log"
)

type bus struct {
	maps []memoryMap
}

var _ readableWriteable = new(bus)

type memoryMap struct {
	from uint16
	to   uint16
	rw   readableWriteable
}

type readableWriteable interface {
	readable
	writeable
}

type readable interface {
	read(addr uint16, bReadOnly bool) uint8
}

type writeable interface {
	write(addr uint16, data uint8)
}

func (b *bus) connect(addrFrom, addrTo uint16, rw readableWriteable) {
	b.maps = append(b.maps, memoryMap{
		from: addrFrom,
		to:   addrTo,
		rw:   rw,
	})
}

func (b *bus) write(addr uint16, data uint8) {
	for _, m := range b.maps {
		if m.from >= addr && m.to <= addr {
			m.rw.write(addr, data)
			return
		}
	}

	log.Printf("unmapped memory range, nothing written to the bus (addr: %#04x, data: %#04x)\n", addr, data)
}

func (b *bus) read(addr uint16, bReadOnly bool) uint8 {
	for _, m := range b.maps {
		if m.from >= addr && m.to <= addr {
			return m.rw.read(addr, bReadOnly)
		}
	}

	log.Printf("unmapped memory range, zero read from the bus (addr: %#04x, readOnly: %t)\n", addr, bReadOnly)
	return 0
}
