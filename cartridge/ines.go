package cartridge

type INESFileHeader struct {
	Magic        uint32  // iNES magic number
	PRGRomChunks uint8   // number of PRG-ROM banks (16KB each)
	CHRRomChunks uint8   // number of CHR-ROM banks (8KB each)
	Control1     uint8   // control bits (mapper, mirroring, battery, trainer)
	Control2     uint8   // control bits (mapper, VS/Playchoice, NES 2.0)
	PRGRamSize   uint8   // PRG-RAM size (8KB each)
	TVSystem1    uint8   // TV system
	TVSystem2    uint8   // TV system
	_            [5]byte // unused
}

func LoadCartridge(path string) (*Cartridge, error) {
	panic("TODO")
}
