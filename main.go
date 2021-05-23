package main

import (
	"fmt"

	"github.com/mtojek/nes-emulator/nes"
)

func main() {
	console := nes.Create()
	fmt.Println(console)
}
