package ui

import (
	"fmt"
	"log"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"

	"gb-player/gb"
)

type Controller struct {
	window    *sdl.Window
	renderer  *sdl.Renderer
	font      *ttf.Font
	view      *View
	timestamp uint64
}

func NewController(window *sdl.Window, renderer *sdl.Renderer, font *ttf.Font) *Controller {
	controller := Controller{}
	controller.window = window
	controller.renderer = renderer
	controller.font = font
	return &controller
}

func (c *Controller) Start(path string) {
	console, err := gb.NewConsole(path)
	if err != nil {
		log.Fatal(err)
	}

	c.view = NewView(c, console)
	c.Run()
}

func (c *Controller) Step() {
	timestamp := sdl.GetTicks64()
	dt := timestamp - c.timestamp
	c.timestamp = timestamp
	c.view.Update(dt)
}

func (c *Controller) Run() {
	var fps float64
	var frameCount uint32
	var frameStart = sdl.GetTicks64()

	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch e := event.(type) {
			case *sdl.QuitEvent:
				running = false
			case *sdl.KeyboardEvent:
				if e.Type == sdl.KEYDOWN {
					if e.Keysym.Sym == sdl.K_ESCAPE {
						running = false
					}
				}
			}
		}

		// Clear screen
		c.renderer.SetDrawColor(0x63, 0x94, 0xED, 0xff)
		c.renderer.Clear()

		// Update state
		c.Step()
		c.drawDebugWindow(fps)

		// Present to screen
		c.renderer.Present()

		// Calculate framerate
		frameCount++
		frameEnd := sdl.GetTicks64()
		elapsed := frameEnd - frameStart

		if elapsed >= 1000 {
			fps = float64(frameCount) / (float64(elapsed) / 1000.0)

			// Reset framerate counters
			frameCount = 0
			frameStart = frameEnd
		}
	}
}

func (c *Controller) drawDebugWindow(fps float64) {
	renderer := c.renderer
	font := c.font

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
