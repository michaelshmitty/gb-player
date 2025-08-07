package gb

import (
	"image"
	"image/color"
)

const (
	ConsoleWidth  = 160
	ConsoleHeight = 144
)

type Console struct {
	Cartridge *Cartridge
	front     *image.RGBA
}

func NewConsole(path string) (*Console, error) {
	cartridge := InsertCartridge(path)
	buffer := image.NewRGBA(image.Rect(0, 0, ConsoleWidth, ConsoleHeight))

	console := Console{Cartridge: cartridge, front: buffer}

	return &console, nil
}

func (console *Console) StepMilliSeconds(ms uint64) {
	speed := int(ms / 3)
	for y := 0; y < ConsoleHeight; y++ {
		for x := 0; x < ConsoleWidth; x++ {
			console.front.Set(x, y, color.RGBA{0, uint8(y + speed), uint8(x + speed), 255})
		}
	}

}

func (console *Console) Buffer() *image.RGBA {

	return console.front
}
