package ppu

type uint8bits uint8

func (b uint8bits) withBit(f uint8, value bool) uint8bits {
	if value {
		return b | uint8bits(f)
	}
	return b & ^uint8bits(f)
}

func (b uint8bits) bit(f uint8) bool {
	val := b & (1 << f)
	return val > 0
}

type uint16bits uint16

func (b uint16bits) withBit(f uint8, value bool) uint16bits {
	if value {
		return b | uint16bits(f)
	}
	return b & ^uint16bits(f)
}

func (b uint16bits) bit(f uint16) bool {
	val := b & (1 << f)
	return val > 0
}
