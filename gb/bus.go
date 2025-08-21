package gb

import (
	"fmt"
)

type Bus struct {
	Cart *Cartridge
}

func NewBus(cart *Cartridge) *Bus {
	bus := Bus{Cart: cart}

	return &bus
}

func (bus *Bus) Read(address uint16) byte {
	if address < 0x8000 {
		// ROM data
		value, err := bus.Cart.Read(address)
		if err != nil {
			fmt.Printf("Error reading from bus address %X: %s", address, err)
			return 0
		}
		return value
	} else {
		fmt.Printf("Reading from bus address %X not implemented!\n", address)
		return 0
	}
}

func (bus *Bus) Write(address uint16, value byte) error {
	if address < 0x8000 {
		// ROM data
		err := bus.Cart.Write(address, value)
		if err != nil {
			return fmt.Errorf("Error writing to bus address %X: %s", address, err)
		}
		return nil
	} else {
		return fmt.Errorf("Writing to bus address %X not implemented!", address)
	}
}
