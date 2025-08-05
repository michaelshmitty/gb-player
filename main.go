package main

import (
	"fmt"
	"log"
	"runtime"

	"gb-player/gb"

	"github.com/AllenDang/cimgui-go/backend"
	"github.com/AllenDang/cimgui-go/backend/sdlbackend"
	"github.com/AllenDang/cimgui-go/imgui"
)

var currentBackend backend.Backend[sdlbackend.SDLWindowFlags]

func init() {
	runtime.LockOSThread()
}

func showDebugWindow() {
	imgui.Text(fmt.Sprintf("Framerate: %.1f FPS", imgui.CurrentIO().Framerate()))
}

func showROMWindow(cartridge *gb.Cartridge) {
	imgui.Begin("ROM")

	imgui.Text(fmt.Sprintf("Title: %s", cartridge.Title))
	imgui.Text(fmt.Sprintf("Mapper: %s", cartridge.Mapper))
	imgui.Text(fmt.Sprintf("Licensee: %s", cartridge.Licensee))
	if cartridge.SGBSupport {
		imgui.Text("Super Game Boy support?: Yes")
	} else {
		imgui.Text("Super Game Boy support?: No")
	}
	imgui.Text(fmt.Sprintf("ROM size: %d bytes", cartridge.ROMSize))
	imgui.Text(fmt.Sprintf("RAM size: %s", cartridge.RAMSize))
	imgui.Text(fmt.Sprintf("Publication: %s", cartridge.Destination))
	imgui.Text(fmt.Sprintf("ROM version: %d", cartridge.Version))

	imgui.End()
}

func showCPUWindow(cpu *gb.Cpu) {
	imgui.Begin("CPU")

	imgui.Text(fmt.Sprintf("PC: %d", cpu.PC))

	imgui.End()
}

func main() {
	currentBackend, err := backend.CreateBackend(sdlbackend.NewSDLBackend())
	if err != nil {
		log.Fatal(err)
	}

	currentBackend.SetBgColor(imgui.NewVec4(0.45, 0.55, 0.6, 1.0))

	currentBackend.CreateWindow("GB Player", 1200, 900)

	currentBackend.SetDropCallback(func(p []string) {
		fmt.Printf("drop triggered: %v", p)
	})

	currentBackend.SetCloseCallback(func() {
		fmt.Println("window is closing")
	})

	cartridge := gb.Insert("rom.gb")
	cpu := gb.Reset()

	currentBackend.Run(func() {
		imgui.ClearSizeCallbackPool()

		gb.Tick(&cpu)

		showDebugWindow()
		showROMWindow(&cartridge)
		showCPUWindow(&cpu)
	})
}
