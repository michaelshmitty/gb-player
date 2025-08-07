package ui

import (
	"log"
	"runtime"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

const (
	windowWidth  = 800
	windowHeight = 600
)

func init() {
	runtime.LockOSThread()
}

func Run(romPath string) {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		log.Fatal(err)
	}
	defer sdl.Quit()

	if err := ttf.Init(); err != nil {
		log.Fatal((err))
	}
	defer ttf.Quit()

	window, err := sdl.CreateWindow(
		"GB Player",
		sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED,
		windowWidth, windowHeight,
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

	font, err := ttf.OpenFont("SourceCodePro.ttf", 18)
	if err != nil {
		log.Fatal(err)
	}
	defer font.Close()
	font.SetStyle(ttf.STYLE_BOLD)

	controller := NewController(window, renderer, font)
	controller.Start(romPath)
}
