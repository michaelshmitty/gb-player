package gb

import (
	"fmt"
	"os"
)

var timer = Timer{DIV: 0xAC00}

type Timer struct {
	DIV  uint16
	TIMA byte
	TMA  byte
	TAC  byte
}

func (timer *Timer) Tick() {
	previousDIV := timer.DIV
	timer.DIV++

	timerNeedsUpdate := false

	// Determine clock mode
	// TAC (Timer control register)
	// Bits 0-1 determines clock frequency
	// 00: 4096 Hz
	// 01: 262144 Hz
	// 10: 65536 Hz
	// 11: 16384 Hz
	// Bit 2 determines whether timer is enabled or not

	clockMode := timer.TAC & (0b11)
	switch clockMode {
	case 0b00:
		// 4096 Hz
		// Detect a falling edge by comparing bit 9 in the previous DIV value
		// with bit 9 in the current DIV value
		previouslyEnabled := previousDIV&(1<<9) != 0
		currentlyDisabled := timer.DIV&(1<<9) == 0
		timerNeedsUpdate = previouslyEnabled && currentlyDisabled

	case 0b01:
		// 262144 Hz
		// Detect a falling edge by comparing bit 3 in the previous DIV value
		// with bit 3 in the current DIV value
		previouslyEnabled := previousDIV&(1<<3) != 0
		currentlyDisabled := timer.DIV&(1<<3) == 0
		timerNeedsUpdate = previouslyEnabled && currentlyDisabled

	case 0b10:
		// 65536 Hz
		// Detect a falling edge by comparing bit 5 in the previous DIV value
		// with bit 5 in the current DIV value
		previouslyEnabled := previousDIV&(1<<5) != 0
		currentlyDisabled := timer.DIV&(1<<5) == 0
		timerNeedsUpdate = previouslyEnabled && currentlyDisabled

	case 0b11:
		// 16384 Hz
		// Detect a falling edge by comparing bit 7 in the previous DIV value
		// with bit 7 in the current DIV value
		previouslyEnabled := previousDIV&(1<<7) != 0
		currentlyDisabled := timer.DIV&(1<<7) == 0
		timerNeedsUpdate = previouslyEnabled && currentlyDisabled
	}

	timerIsEnabled := timer.TAC&(1<<2) != 0
	// If the timer needs to be updated based on the determined clock mode and the timer is enabled, increment the timer
	if timerNeedsUpdate && timerIsEnabled {
		timer.TIMA++

		// Check if TIMA is going to wrap and trigger an interrupt if necessary
		if timer.TIMA == 0xFF {
			timer.TIMA = timer.TMA

			fmt.Println("TODO: cpu_request_interrupt(IT_TIMER")
		}
	}

}

func (timer *Timer) Read(address uint16) byte {
	switch address {
	case 0xFF04:
		return byte(timer.DIV >> 8)

	case 0xFF05:
		return timer.TIMA

	case 0xFF06:
		return timer.TMA

	case 0xFF07:
		return timer.TAC

	default:
		fmt.Printf("Reading from Timer: invalid address %X\n", address)
		os.Exit(1)
	}
	return 0
}

func (timer *Timer) Write(address uint16, value byte) {
	switch address {
	case 0xFF04:
		//DIV
		timer.DIV = 0

	case 0xFF05:
		//TIMA
		timer.TIMA = value

	case 0xFF06:
		//TMA
		timer.TMA = value

	case 0xFF07:
		//TAC
		timer.TAC = value

	default:
		fmt.Printf("Writing to Timer: invalid address %X\n", address)
		os.Exit(1)
	}
}
