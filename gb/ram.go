package gb

import (
	"fmt"
	"os"
)

type RAM struct {
	WRAM [0x2000]byte
	HRAM [0x80]byte
}

func NewRAM() *RAM {
	ram := RAM{}

	return &ram
}

func (ram *RAM) WRAMRead(address uint16) byte {
	// TODO(m): Understand this line
	address -= 0xC000

	if address >= 0x2000 {
		fmt.Printf("Reading from WRAM: invalid address %X\n", address+0xC000)
		os.Exit(1)
	}

	return ram.WRAM[address]
}

func (ram *RAM) WRAMWrite(address uint16, value byte) {
	// TODO(m): Understand this line
	address -= 0xC000

	if address >= 0x2000 {
		fmt.Printf("Writing to WRAM: invalid address %X\n", address)
		os.Exit(1)
	}

	ram.WRAM[address] = value
}

func (ram *RAM) HRAMRead(address uint16) byte {
	// TODO(m): Understand this line
	address -= 0xFF80

	return ram.HRAM[address]
}

func (ram *RAM) HRAMWrite(address uint16, value byte) {
	// TODO(m): Understand this line
	address -= 0xFF80

	ram.HRAM[address] = value
}
