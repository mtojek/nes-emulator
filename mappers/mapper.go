package mappers

import (
	"fmt"

	"github.com/mtojek/nes-emulator/bus"
)

type Mapper interface {
	ID() uint8

	ConnectTo(cpuBus *bus.Bus, ppuBus *bus.Bus, prgMemory bus.ReadableWriteable, chrMemory bus.ReadableWriteable)
	MapCPU(addr uint16) uint16
	MapPPU(addr uint16) uint16
}

func Load(mapperID uint8, nPRGBanks uint8) (Mapper, error) {
	if mapperID == 0 {
		return &mapper000{
			nPRGBanks: nPRGBanks,
		}, nil
	}
	return nil, fmt.Errorf("unsupported mapper: %d", mapperID)
}
