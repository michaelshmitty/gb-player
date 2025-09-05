package gb

import (
	"fmt"
	"os"
)

var SerialData [2]byte
var InterruptFlags byte

func IORead(address uint16) byte {
	if address == 0xFF01 {
		return SerialData[0]
	} else if address == 0xFF02 {
		return SerialData[1]
	} else if (address >= 0xFF04) && (address <= 0xFF07) {
		return timer.Read(address)
	} else if address == 0xFF0F {
		return InterruptFlags
	} else if (address >= 0xFF10) && (address <= 0xFF3F) {
		// Ignore sound
		return 0
	} else {
		fmt.Printf("Reading from IO: invalid address %X\n", address)
		os.Exit(1)
	}
	return 0
}

func IOWrite(address uint16, value byte) {
	if address == 0xFF01 {
		SerialData[0] = value
	} else if address == 0xFF02 {
		SerialData[1] = value
	} else if (address >= 0xFF04) && (address <= 0xFF07) {
		timer.Write(address, value)
	} else if address == 0xFF0F {
		InterruptFlags = value
	} else if (address >= 0xFF10) && (address <= 0xFF3F) {
		// Ignore sound
	} else {
		fmt.Printf("Writing to IO: invalid address %X\n", address)
		os.Exit(1)
	}
}
