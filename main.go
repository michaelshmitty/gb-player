package main

import (
	"fmt"
	"log"
	"runtime"

	"gb-player/gb"

	"github.com/veandco/go-sdl2/sdl"
)

func init() {
	runtime.LockOSThread()
}

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
	fmt.Println()

	cpu := gb.Reset()
	fmt.Println("Program counter:", cpu.PC)
	gb.Tick(&cpu)
	fmt.Println("Program counter:", cpu.PC)

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		log.Fatal(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow(
		"GB Player",
		sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED,
		800, 600,
		sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, 0)
	if err != nil {
		panic(err)
	}
	defer renderer.Destroy()

	renderer.RenderSetVSync(true)
	renderer.SetDrawColor(0x63, 0x94, 0xED, 0xff)

	var frameCount uint32 = 0
	var startTicks = sdl.GetTicks()

	running := true
	for running {
		// frameStartTime := sdl.GetTicks()

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				running = false
				break
			}
		}

		renderer.Clear()
		renderer.Present()

		// Calculate framerate
		frameCount++
		currentTicks := sdl.GetTicks()
		elapsed := currentTicks - startTicks
		if elapsed >= 1000 {
			fps := float32(frameCount) / (float32(elapsed) / 1000.0)
            fmt.Printf("FPS: %.2f\n", fps)

            // Reset counters
            frameCount = 0
            startTicks = currentTicks
		}
	}
}
