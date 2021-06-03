package controller

import "github.com/mtojek/nes-emulator/bus"

const (
	ButtonA = iota
	ButtonB
	ButtonSelect
	ButtonStart
	ButtonUp
	ButtonDown
	ButtonLeft
	ButtonRight
)

type Controller struct {
	buttons [8]bool
	index   byte
	strobe  byte
}

var _ bus.ReadableWriteable = new(Controller)

func NewController() *Controller {
	return new(Controller)
}

func (c *Controller) Read(_ uint16) uint8 {
	value := byte(0)
	if c.index < 8 && c.buttons[c.index] {
		value = 1
	}
	c.index++
	if c.strobe&1 == 1 {
		c.index = 0
	}
	return value
}

func (c *Controller) Write(_ uint16, data uint8) {
	c.strobe = data
	if c.strobe&1 == 1 {
		c.index = 0
	}
}

func (c *Controller) SetState(buttons [8]bool) {
	c.buttons = buttons
}
