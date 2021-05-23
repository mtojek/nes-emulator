package cartridge

import (
	"encoding/binary"
	"io"
	"os"

	"github.com/mtojek/nes-emulator/bus"
	"github.com/mtojek/nes-emulator/mappers"
	"github.com/pkg/errors"
)

const (
	HORIZONTAL = iota
	VERTICAL
	ONESCREEN_LO
	ONESCREEN_HI
)

type Cartridge struct {
	vPRGMemory []uint8
	vCHRMemory []uint8

	mapper mappers.Mapper
	mirror uint8

	nPRGBanks uint8
	nCHRBanks uint8
}

var _ bus.ReadableWriteable = new(Cartridge)

func Load(path string) (*Cartridge, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrap(err, "can't open file")
	}
	defer file.Close()

	var header INESFileHeader
	if err := binary.Read(file, binary.LittleEndian, &header); err != nil {
		return nil, err
	}

	// verify header magic number
	if header.Magic != iNESFileMagic {
		return nil, errors.New("invalid iNES file (bad magic)")
	}

	// read trainer if present (unused)
	if header.Control1&4 == 4 {
		trainer := make([]byte, 512)
		if _, err := io.ReadFull(file, trainer); err != nil {
			return nil, errors.Wrap(err, "can't read trainer")
		}
	}

	// Read PRG-ROM bank(s)
	prg := make([]byte, int(header.PRGROMChunks)*0x4000)
	if _, err := io.ReadFull(file, prg); err != nil {
		return nil, errors.Wrap(err, "can't read PRG-ROM banks")
	}

	// Read CHR-ROM bank(s)
	chr := make([]byte, int(header.CHRROMChunks)*0x2000)
	if _, err := io.ReadFull(file, chr); err != nil {
		return nil, errors.Wrap(err, "can't read CHR-ROM banks")
	}

	// Configure mapper
	nMapperID := ((header.Control2 >> 4) << 4) | (header.Control1 >> 4)
	mapper, err := mappers.Load(nMapperID)
	if err != nil {
		return nil, errors.Wrap(err, "can't load selected mapper")
	}

	// Configure mirroring
	mirror1 := header.Control1 & 1
	mirror2 := (header.Control1 >> 3) & 1
	mirror := mirror1 | mirror2<<1

	return &Cartridge{
		vPRGMemory: prg,
		vCHRMemory: chr,

		mapper: mapper,
		mirror: mirror,

		nPRGBanks: header.PRGROMChunks,
		nCHRBanks: header.CHRROMChunks,
	}, nil
}

func (c *Cartridge) Read(addr uint16, bReadOnly bool) uint8 {
	mapped := c.mapper.Map(addr)
	return c.vPRGMemory[mapped]
}

func (c *Cartridge) Write(addr uint16, data uint8) {
	mapped := c.mapper.Map(addr)
	c.vCHRMemory[mapped] = data
}

func (c *Cartridge) ConnectTo(cpuBus *bus.Bus, ppuBus *bus.Bus) {
	c.mapper.ConnectTo(cpuBus, ppuBus, c)
}
