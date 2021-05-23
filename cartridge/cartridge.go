package cartridge

import (
	"encoding/binary"
	"github.com/pkg/errors"
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
	nesFile, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrap(err, "can't open file")
	}

	var header INESFileHeader
	if err := binary.Read(nesFile, binary.LittleEndian, &header); err != nil {
		return nil, err
	}

	// verify header magic number
	if header.Magic != iNESFileMagic {
		return nil, errors.New("invalid iNES file (bad magic)")
	}

	return &Cartridge{}, nil
}
