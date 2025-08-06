package main

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

func main() {
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

	var fps float64
	var frameCount uint32
	var startTicks = sdl.GetTicks64()
	var now uint64
	var last uint64

	cartridge := Insert("rom.gb")
	fmt.Println(cartridge)

	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				running = false
			}
		}

		last = now
		now = sdl.GetPerformanceCounter()
		deltaTime := float64((now - last) * 1000.0 / sdl.GetPerformanceFrequency())
		_ = deltaTime

		// Clear screen
		renderer.SetDrawColor(0x63, 0x94, 0xED, 0xff)
		renderer.Clear()

		// Update state
		rect := sdl.Rect{X: 100, Y: 100, W: 200, H: 150}
		renderer.SetDrawColor(0, 0, 255, 255)
		renderer.FillRect(&rect)

		drawDebugWindow(fps)

		// Present to screen
		renderer.Present()

		// Calculate framerate
		frameCount++
		currentTicks := sdl.GetTicks64()
		elapsed := currentTicks - startTicks

		if elapsed >= 1000 {
			fps = float64(frameCount) / (float64(elapsed) / 1000.0)

			// Reset counters
			frameCount = 0
			startTicks = currentTicks
		}
	}
}

func drawDebugWindow(fps float64) {
	var textSurface *sdl.Surface
	var textTexture *sdl.Texture

	// FPS
	textSurface, err := font.RenderUTF8Blended(fmt.Sprintf("FPS: %.2f", fps), sdl.Color{R: 0, G: 0, B: 0, A: 255})
	if err != nil {
		log.Fatal(err)
	}
	defer textSurface.Free()

	textTexture, err = renderer.CreateTextureFromSurface(textSurface)
	if err != nil {
		log.Fatal(err)
	}
	defer textTexture.Destroy()

	textRect := sdl.Rect{X: 0, Y: 0, W: textSurface.W, H: textSurface.H}
	renderer.Copy(textTexture, nil, &textRect)
}
