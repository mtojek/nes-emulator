package cartridge

import (
	"encoding/binary"
	"github.com/pkg/errors"
	"io"
	"os"
)

type Cartridge struct {
	vPRGMemory []uint8
	vCHRMemory []uint8

	nMapperID uint8
	nPRGBanks uint8
	nCHRBanks uint8
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

	// Read PRG-ROM bank(s)
	prg := make([]byte, int(header.PRGRomChunks)*0x4000)
	if _, err := io.ReadFull(file, prg); err != nil {
		return nil, errors.Wrap(err, "can't read PRG-ROM banks")
	}

	// Read CHR-ROM bank(s)
	chr := make([]byte, int(header.CHRRomChunks)*0x2000)
	if _, err := io.ReadFull(file, chr); err != nil {
		return nil, errors.Wrap(err, "can't read CHR-ROM banks")
	}

	nMapperID := ((header.Control2 >> 4) << 4) | (header.Control1 >> 4)

	return &Cartridge{
		vPRGMemory: prg,
		vCHRMemory: chr,

		nMapperID: nMapperID,
		nPRGBanks: header.PRGRomChunks,
		nCHRBanks: header.CHRRomChunks,
	}, nil
}
