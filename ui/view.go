package ui

import (
	"image"
	"log"

	"github.com/veandco/go-sdl2/sdl"

	"gb-player/gb"
)

type View struct {
	controller *Controller
	console    *gb.Console
}

func NewView(controller *Controller, console *gb.Console) *View {
	return &View{controller, console}
}

func (view *View) Update(dt uint64) {
	console := view.console
	renderer := view.controller.renderer

	console.StepMilliSeconds(dt)

	buffer := console.Buffer()

	drawBuffer(renderer, buffer)
}

func drawBuffer(renderer *sdl.Renderer, buffer *image.RGBA) {
	width, height := gb.ConsoleWidth, gb.ConsoleHeight

	imageTexture, err := renderer.CreateTexture(
		uint32(sdl.PIXELFORMAT_RGBA32),
		sdl.TEXTUREACCESS_STREAMING,
		int32(width), int32(height))
	if err != nil {
		log.Fatal(err)
	}
	defer imageTexture.Destroy()

	pixels, pitch, err := imageTexture.Lock(nil)
	if err != nil {
		log.Fatal(err)
	}
	for y := range height {
		start := y * buffer.Stride
		end := start + width*4
		copy(pixels[y*pitch:y*pitch+width*4], buffer.Pix[start:end])
	}
	imageTexture.Unlock()

	x := (windowWidth / 2) - width
	y := (windowHeight / 2) - height
	imageRect := sdl.Rect{X: int32(x), Y: int32(y), W: 2 * int32(width), H: 2 * int32(height)}

	renderer.Copy(imageTexture, nil, &imageRect)
}
