package main

import (
	"gb-player/ui"
)

func main() {
	// FIXME(m): Allow specifying rom file on command line
	// if len(os.Args) != 2 {
	// 	log.Fatalln("No rom file specified")
	// }
	// romPath := os.Args[1]
	romPath := "./roms/dmg-acid2.gb"

	ui.Run(romPath)
}
