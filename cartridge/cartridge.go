package cartridge

import (
	"encoding/binary"
	"io"
	"os"

	"github.com/mtojek/nes-emulator/bus"
	"github.com/mtojek/nes-emulator/mappers"
	"github.com/pkg/errors"
)

type Cartridge struct {
	vPRGMemory []uint8
	vCHRMemory []uint8

	mapper mappers.Mapper
	mirror uint8

	nPRGBanks uint8
	nCHRBanks uint8
}

type prgMemoryHandler struct {
	c *Cartridge
}

var _ bus.ReadableWriteable = new(prgMemoryHandler)

type chrMemoryHandler struct {
	c *Cartridge
}

var _ bus.ReadableWriteable = new(chrMemoryHandler)

func (mh *prgMemoryHandler) Read(addr uint16) uint8 {
	mapped := mh.c.mapper.MapCPU(addr)
	return mh.c.vPRGMemory[mapped]
}

func (mh *prgMemoryHandler) Write(addr uint16, data uint8) {
	mapped := mh.c.mapper.MapCPU(addr)
	mh.c.vPRGMemory[mapped] = data
}

func (mh *chrMemoryHandler) Read(addr uint16) uint8 {
	mapped := mh.c.mapper.MapPPU(addr)
	return mh.c.vCHRMemory[mapped]
}

func (mh *chrMemoryHandler) Write(addr uint16, data uint8) {
	mapped := mh.c.mapper.MapPPU(addr)
	mh.c.vCHRMemory[mapped] = data
}

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

	// Read prg-ROM bank(s)
	prg := make([]byte, int(header.PRGROMChunks)*0x4000)
	if _, err := io.ReadFull(file, prg); err != nil {
		return nil, errors.Wrap(err, "can't read prg-ROM banks")
	}

	// Read CHR-ROM bank(s)
	chr := make([]byte, int(header.CHRROMChunks)*0x2000)
	if _, err := io.ReadFull(file, chr); err != nil {
		return nil, errors.Wrap(err, "can't read CHR-ROM banks")
	}

	// Configure mapper
	nMapperID := ((header.Control2 >> 4) << 4) | (header.Control1 >> 4)
	mapper, err := mappers.Load(nMapperID, header.PRGROMChunks)
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

func (c *Cartridge) prg() bus.ReadableWriteable {
	return &prgMemoryHandler{c}
}

func (c *Cartridge) chr() bus.ReadableWriteable {
	return &chrMemoryHandler{c}
}

func (c *Cartridge) ConnectTo(cpuBus *bus.Bus, ppuBus *bus.Bus) {
	c.mapper.ConnectTo(cpuBus, ppuBus, c.prg(), c.chr())
}

func (c *Cartridge) Mirroring() uint8 {
	return c.mirror
}