package main

import (
	"fmt"
	"gb-player/gb"
)

func main() {
	cartridge := gb.Insert("rom.gb")
	fmt.Println("Title:", cartridge.Title)
	fmt.Println("Mapper:", cartridge.Mapper)
	fmt.Println("Licensee:", cartridge.Licensee)
	fmt.Println("Super Game Boy support?", cartridge.SGBSupport)
	fmt.Println("ROM size:", cartridge.ROMSize, "bytes")
	fmt.Println("RAM size:", cartridge.RAMSize)
	fmt.Println("Publication destination:", cartridge.Destination)
	fmt.Println("ROM version:", cartridge.Version)
}
