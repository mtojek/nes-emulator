package mappers

import (
	"fmt"

	"github.com/mtojek/nes-emulator/bus"
)

type Mapper interface {
	ID() uint8

	ConnectTo(cpuBus *bus.Bus, ppuBus *bus.Bus, rw bus.ReadableWriteable)
	Map(addr uint16) uint16
}

func Load(mapperID uint8) (Mapper, error) {
	if mapperID == 0 {
		return new(mapper000), nil
	}
	return nil, fmt.Errorf("unsupported mapper: %d", mapperID)
}
