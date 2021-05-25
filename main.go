package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/mtojek/nes-emulator/cartridge"
	"github.com/mtojek/nes-emulator/nes"
	"github.com/mtojek/nes-emulator/ui"
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

	window, tex, err := ui.Initialize()
	if err != nil {
		log.Fatal(err)
	}
	defer ui.Terminate(window)

	for !window.ShouldClose() {
		console.DrawNewFrame()

		startFrameTime := time.Now()
		// PPU processing
		for !console.FrameComplete() {
			console.Clock()
		}

		// OpenGL processing
		ui.Redraw(window, tex, console.Buffer())

		processingDuration := time.Now().Sub(startFrameTime)
		waitingTime := time.Second/60 - processingDuration

		if waitingTime > 0 {
			fmt.Printf("Sleep for: %v\n", waitingTime)
			time.Sleep(waitingTime)
		}
	}
}
