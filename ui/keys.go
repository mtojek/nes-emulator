package ui

import (
	"github.com/mtojek/nes-emulator/controller"

	"github.com/go-gl/glfw/v3.2/glfw"
)

func ReadKeysPlayer1(window *glfw.Window) [8]bool {
	var result [8]bool
	result[controller.ButtonB] = readKey(window, glfw.KeyZ)
	result[controller.ButtonA] = readKey(window, glfw.KeyX)
	result[controller.ButtonSelect] = readKey(window, glfw.KeyRightShift)
	result[controller.ButtonStart] = readKey(window, glfw.KeyEnter)
	result[controller.ButtonUp] = readKey(window, glfw.KeyUp)
	result[controller.ButtonDown] = readKey(window, glfw.KeyDown)
	result[controller.ButtonLeft] = readKey(window, glfw.KeyLeft)
	result[controller.ButtonRight] = readKey(window, glfw.KeyRight)
	return result
}

func readKey(window *glfw.Window, key glfw.Key) bool {
	return window.GetKey(key) == glfw.Press
}
