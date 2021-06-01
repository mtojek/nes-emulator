package bus

import (
	"fmt"
	"log"
)

type Bus struct {
	maps []memoryMap
}

var _ ReadableWriteable = new(Bus)

type memoryMap struct {
	from uint16
	to   uint16
	rw   ReadableWriteable
}

type ReadableWriteable interface {
	Readable
	Writeable
}

type Readable interface {
	Read(addr uint16) uint8
}

type Writeable interface {
	Write(addr uint16, data uint8)
}

func (b *Bus) Connect(addrFrom, addrTo uint16, rw ReadableWriteable) {
	b.maps = append(b.maps, memoryMap{
		from: addrFrom,
		to:   addrTo,
		rw:   rw,
	})
}

func (b *Bus) Read(addr uint16) uint8 {
	for _, m := range b.maps {
		if m.from <= addr && addr <= m.to {
			return m.rw.Read(addr)
		}
	}

	log.Printf("unmapped memory range, zero read from the Bus (addr: %#04x)\n", addr)
	return 0
}

func (b *Bus) Write(addr uint16, data uint8) {
	for _, m := range b.maps {
		if m.from <= addr && addr <= m.to {
			m.rw.Write(addr, data)
			return
		}
	}

	log.Printf("unmapped memory range, nothing written to the Bus (addr: %#04x, data: %#02x)\n", addr, data)
}

func (b *Bus) Print(from, to uint16) {
	for offset := from; offset <= to; offset++ {
		if offset%16 == 0 {
			fmt.Printf("%04X: ", offset)
		}
		fmt.Printf("%02X ", b.Read(offset))

		if offset%16 == 15 {
			fmt.Println()
		}

		if offset == 0xFFFF {
			break // safety exit
		}
	}
}
