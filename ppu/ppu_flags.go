package ppu

const (
	flagStatusUnused         = uint8(1) << 0
	flagStatusSpriteOverflow = uint8(1) << 5
	flagStatusSpriteZeroHit  = uint8(1) << 6
	flagStatusVerticalBlank  = uint8(1) << 7
)

const (
	flagMaskGrayscale            = uint8(1) << 0
	flagMaskRenderBackgroundLeft = uint8(1) << 1
	flagMaskRenderSpritesLeft    = uint8(1) << 2
	flagMaskRenderBackground     = uint8(1) << 3
	flagMaskRenderSprites        = uint8(1) << 4
	flagMaskEnhanceRed           = uint8(1) << 5
	flagMaskEnhanceGreen         = uint8(1) << 6
	flagMaskEnhanceBlue          = uint8(1) << 7
)

const (
	flagControlNametableX        = uint8(1) << 0
	flagControlNametableY        = uint8(1) << 1
	flagControlIncrementMode     = uint8(1) << 2
	flagControlPatternSprite     = uint8(1) << 3
	flagControlPatternBackground = uint8(1) << 4
	flagControlSpriteSize        = uint8(1) << 5
	flagControlSlaveMode         = uint8(1) << 6
	flagControlEnableNMI         = uint8(1) << 7
)

const (
	flagLoopyCoarseX = uint8(1) << 0
	flagLoopyCoarseY = uint8(1) << 5
	flagLoopyNametableX = uint8(1) << 10
	flagLoopyNametableY = uint8(1) << 11
	flagLoopyFineY = uint8(1) << 12
	flagLoopyUnused = uint8(1) << 15
)