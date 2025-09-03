package gb

import (
	"fmt"
	"os"
)

var SerialData [2]byte

func IORead(address uint16) byte {
	if address == 0xFF01 {
		return SerialData[0]
	}

	if address == 0xFF02 {
		return SerialData[1]
	}

	if (address >= 0xFF04) && (address <= 0xFF07) {
		fmt.Printf("Reading from IO: invalid address %X. Timer not yet implemented!", address)
		os.Exit(1)
	}

	fmt.Printf("Reading from IO: invalid address %X", address)
	os.Exit(1)

	return 0
}

func IOWrite(address uint16, value byte) {
	if address == 0xFF01 {
		SerialData[0] = value
	}

	if address == 0xFF02 {
		SerialData[1] = value
	}

	if (address >= 0xFF04) && (address <= 0xFF07) {
		fmt.Printf("Writing to IO: invalid address %X. Timer not yet implemented!", address)
		os.Exit(1)
	}

	fmt.Printf("Writing to IO: invalid address %X", address)
	os.Exit(1)
}
