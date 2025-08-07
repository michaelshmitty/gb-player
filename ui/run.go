package ui

import (
	"fmt"
	"log"
	"runtime"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

var renderer *sdl.Renderer
var font *ttf.Font

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
		800, 600,
		sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	renderer, err = sdl.CreateRenderer(window, -1, 0)
	if err != nil {
		panic(err)
	}
	defer renderer.Destroy()
	renderer.RenderSetVSync(true)

	font, err = ttf.OpenFont("SourceCodePro.ttf", 18)
	if err != nil {
		log.Fatal(err)
	}
	defer font.Close()
	font.SetStyle(ttf.STYLE_BOLD)

	controller := NewController(renderer)
	controller.Start(romPath)
}

func drawDebugWindow(fps float64) {
	// FPS
	textSurface, err := font.RenderUTF8Blended(fmt.Sprintf("FPS: %.2f", fps), sdl.Color{R: 0, G: 0, B: 0, A: 255})
	if err != nil {
		log.Fatal(err)
	}
	defer textSurface.Free()

	textTexture, err := renderer.CreateTextureFromSurface(textSurface)
	if err != nil {
		log.Fatal(err)
	}
	defer textTexture.Destroy()

	textRect := sdl.Rect{X: 0, Y: 0, W: textSurface.W, H: textSurface.H}
	renderer.Copy(textTexture, nil, &textRect)
}
