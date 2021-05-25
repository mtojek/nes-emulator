package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/mtojek/nes-emulator/cartridge"
	"github.com/mtojek/nes-emulator/nes"
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

	for {
		console.DrawNewFrame()

		startFrameTime := time.Now()
		for console.FrameComplete() {
			console.Clock()
		}
		drawingDuration := time.Now().Sub(startFrameTime)
		waitingTime := time.Second / 60 - drawingDuration

		if waitingTime > 0 {
			time.Sleep(waitingTime)
		}
	}
}
