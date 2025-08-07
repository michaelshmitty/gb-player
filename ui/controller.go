package ui

import (
	"gb-player/gb"

	"github.com/veandco/go-sdl2/sdl"
)

type Controller struct {
	renderer  *sdl.Renderer
	view      *View
	timestamp uint64
}

func NewController(renderer *sdl.Renderer) *Controller {
	controller := Controller{}
	controller.renderer = renderer
	return &controller
}

func (c *Controller) Start(path string) {
	console := gb.NewConsole(path)
	c.view = NewView(c, console)
	c.Run()
}

func (c *Controller) Step() {
	timestamp := sdl.GetTicks64()
	dt := timestamp - c.timestamp
	c.view.Update(timestamp, dt)
}

func (c *Controller) Run() {
	var fps float64
	var frameCount uint32
	var startTicks = sdl.GetTicks64()
	var now, last uint64

	// cartridge := gb.Insert(romPath)
	// fmt.Println(cartridge)

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
		c.Step()
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
