package cartridge

import (
	"encoding/binary"
	"github.com/mtojek/nes-emulator/bus"
	"github.com/pkg/errors"
	"io"
	"os"
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

	nMapperID uint8
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

	// Configure mirroring
	mirror1 := header.Control1 & 1
	mirror2 := (header.Control1 >> 3) & 1
	mirror := mirror1 | mirror2<<1

	return &Cartridge{
		vPRGMemory: prg,
		vCHRMemory: chr,

		nMapperID: nMapperID,
		mirror: mirror,

		nPRGBanks: header.PRGROMChunks,
		nCHRBanks: header.CHRROMChunks,
	}, nil
}

func (c *Cartridge) Read(addr uint16, bReadOnly bool) uint8 {
	panic("implement me")
}

func (c *Cartridge) Write(addr uint16, data uint8) {
	panic("implement me")
}