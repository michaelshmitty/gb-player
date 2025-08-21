package main

import (
	"fmt"
	"log"
	"os"

	// "gb-player/ui"
	"gb-player/gb"
)

var running = true

func main() {
	// FIXME(m): Allow specifying rom file on command line
	if len(os.Args) != 2 {
		log.Fatalln("No rom file specified")
	}
	romPath := os.Args[1]

	// ui.Run(romPath)

	console, err := gb.NewConsole(romPath)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Executing instructions")
	running := true
	for running {
		console.CPU.Step()
	}
}
