package gb

import (
	"fmt"
	"image"
	"image/color"
)

const (
	ConsoleWidth  = 160
	ConsoleHeight = 144
)

type Console struct {
	Bus   *Bus
	CPU   *CPU
	front *image.RGBA
	Ticks uint64
}

func NewConsole(path string) (*Console, error) {
	cart, err := InsertCartridge(path)
	if err != nil {
		return &Console{}, err
	}

	fmt.Println("Cartridge loaded:")
	fmt.Printf("\t Title    : %s\n", cart.Title)
	fmt.Printf("\t Type     : %02X (%s)\n", cart.Header.CartridgeType, cart.Mapper)
	fmt.Printf("\t ROM Size : %d KB\n", cart.ROMSize)
	fmt.Printf("\t RAM Size : %02X (%s)\n", cart.Header.RAMSize, cart.RAMSize)
	fmt.Printf("\t LIC Code : %02X (%s)\n", cart.Header.OldLicenseeCode, cart.Licensee)
	fmt.Printf("\t ROM Vers : %02X\n", cart.Header.MaskROMVersionNumber)
	fmt.Printf("\t Checksum : %02X ", cart.Checksum)
	if cart.Checksum == cart.Header.Checksum {
		fmt.Printf("(PASSED)\n")
	} else {
		fmt.Printf("(FAILED)\n")
	}

	buffer := image.NewRGBA(image.Rect(0, 0, ConsoleWidth, ConsoleHeight))

	bus := NewBus(cart)
	console := Console{Bus: bus, CPU: NewCPU(bus), front: buffer}

	return &console, nil
}

func (console *Console) StepMilliSeconds(ms uint64) {
	speed := int(ms / 3)
	for y := range ConsoleHeight {
		for x := range ConsoleWidth {
			console.front.Set(x, y, color.RGBA{0, uint8(y + speed), uint8(x + speed), 255})
		}
	}
}

func (console *Console) Buffer() *image.RGBA {

	return console.front
}

func (console *Console) Cycle(cycles int) {
	for range cycles {
		// Cycles are given in M-cycles (machine cycles) instead of T-states (system clock ticks)
		// One machine cycle takes 4 system clock ticks to complete
		for range 4 {
			console.Ticks++
			timer.Tick()
		}
	}
}
