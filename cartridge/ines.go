package cartridge

type INESFileHeader struct {
	Magic        uint32  // iNES magic number
	PRGRomChunks uint8   // number of PRG-ROM banks (16KB each)
	CHRRomChunks uint8   // number of CHR-ROM banks (8KB each)
	Control1     uint8   // control bits
	Control2     uint8   // control bits
	PRGRamSize   uint8   // PRG-RAM size (x 8KB)
	TVSystem1    uint8   // TV system
	TVSystem2    uint8   // TV system
	_            [5]byte // unused
}

func LoadCartridge(path string) (*Cartridge, error) {
	panic("TODO")
}
