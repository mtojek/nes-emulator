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

var keysPlayer2 = [8]bool{false, false, false, false, false, false, false, false} // TODO: unimplemented

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

	sound, err := ui.OpenAudioStream(console.AudioBuffer())
	if err != nil {
		log.Fatal(err)
	}
	defer sound.Close()

	window, tex, err := ui.Initialize()
	if err != nil {
		log.Fatal(err)
	}
	defer ui.Terminate(window)

	for !window.ShouldClose() {
		startFrameTime := time.Now()

		console.DrawNewFrame()
		// PPU processing
		for !console.FrameComplete() {
			console.Clock()
		}

		// Read controller keys
		keysPlayer1 := ui.ReadKeysPlayer1(window)
		console.UpdateControllers(keysPlayer1, keysPlayer2)

		// OpenGL processing
		ui.Redraw(window, tex, console.Buffer())

		processingDuration := time.Now().Sub(startFrameTime)
		waitingTime := time.Second/60 - processingDuration

		if waitingTime > 0 {
			//fmt.Printf("Sleep for: %v\n", waitingTime)
			time.Sleep(waitingTime)
		}
	}
}
