package main

import (
	"flag"
	"fmt"
	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/mtojek/nes-emulator/cartridge"
	"github.com/mtojek/nes-emulator/nes"
	"github.com/pkg/errors"
	"log"
	"os"
	"runtime"
	"time"
)

func main() {
	flag.Parse()

	if flag.NArg() == 0 {
		fmt.Println("usage: nes-emulator game.nes")
		os.Exit(1)
	}

	cart, err := cartridge.Load(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}

	console := nes.Create()
	console.Insert(cart)
	console.Reset()

	runtime.LockOSThread()

	window, err := initGlfw()
	if err != nil {
		log.Fatal(err)
	}
	defer glfw.Terminate()

	program, err := initOpenGL()
	if err != nil {
		log.Fatal(err)
	}

	for !window.ShouldClose() {
		console.DrawNewFrame()

		startFrameTime := time.Now()
		// PPU processing
		for !console.FrameComplete() {
			console.Clock()
		}

		// OpenGL processing
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.UseProgram(program)

		glfw.PollEvents()
		window.SwapBuffers()

		drawingDuration := time.Now().Sub(startFrameTime)
		waitingTime := time.Second/60 - drawingDuration

		if waitingTime > 0 {
			//fmt.Printf("Sleep for: %v\n", waitingTime)
			time.Sleep(waitingTime)
		}
	}
}

func initGlfw() (*glfw.Window, error) {
	if err := glfw.Init(); err != nil {
		return nil, errors.Wrap(err, "can't init GL")
	}

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4) // OR 2
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(256, 240, "NES Emulator", nil, nil)
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