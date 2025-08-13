package main

import (
	"fmt"
	"log"

	// "gb-player/ui"
	"gb-player/gb"
)

var running = true

func main() {
	// FIXME(m): Allow specifying rom file on command line
	// if len(os.Args) != 2 {
	// 	log.Fatalln("No rom file specified")
	// }
	// romPath := os.Args[1]
	romPath := "./roms/dmg-acid2.gb"

	// ui.Run(romPath)

	console, err := gb.NewConsole(romPath)
	if err != nil {
		log.Fatal(err)
	}

	running := true
	for running {
		err := console.CPU.Step()
		if err != nil {
			fmt.Println(err)
			running = false
			fmt.Println("CPU stopped")
		} else {
			fmt.Println("CPU step")
		}
	}
}
