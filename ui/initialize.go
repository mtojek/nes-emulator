package ui

import (
	"image"
	"log"
	"runtime"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/pkg/errors"
)

const (
	windowSizeX = 768
	windowSizeY = 720
)

func init() {
	runtime.LockOSThread()
}

func Initialize() (*glfw.Window, uint32, error) {
	window, err := initGlfw()
	if err != nil {
		return nil, 0, errors.Wrap(err, "initGlfw failed")
	}

	program, err := initOpenGL()
	if err != nil {
		return nil, 0, errors.Wrap(err, "initOpenGL failed")
	}

	tex := createTexture()
	gl.Clear(gl.COLOR_BUFFER_BIT)
	gl.UseProgram(program)

	return window, tex, nil
}

func Redraw(window *glfw.Window, tex uint32, buffer *image.RGBA) {
	gl.BindTexture(gl.TEXTURE_2D, tex)
	setTexture(buffer)
	drawBuffer(window)
	gl.BindTexture(gl.TEXTURE_2D, 0)

	window.SwapBuffers()
	glfw.PollEvents()
}

func Terminate(window *glfw.Window) {
	glfw.Terminate()
}

func initGlfw() (*glfw.Window, error) {
	if err := glfw.Init(); err != nil {
		return nil, errors.Wrap(err, "can't init GL")
	}
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)

	window, err := glfw.CreateWindow(windowSizeX, windowSizeY, "NES Emulator", nil, nil)
	if err != nil {
		return nil, errors.Wrap(err, "can't create GL window")
	}
	window.MakeContextCurrent()
	return window, nil
}

func initOpenGL() (uint32, error) {
	if err := gl.Init(); err != nil {
		return 0, err
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL version", version)

	prog := gl.CreateProgram()
	gl.LinkProgram(prog)
	return prog, nil
}
