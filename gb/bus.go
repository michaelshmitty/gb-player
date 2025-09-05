package gb

import (
	"fmt"
)

// 0x0000 - 0x3FFF : ROM Bank 0
// 0x4000 - 0x7FFF : ROM Bank 1 - Switchable
// 0x8000 - 0x97FF : CHR RAM
// 0x9800 - 0x9BFF : BG Map 1
// 0x9C00 - 0x9FFF : BG Map 2
// 0xA000 - 0xBFFF : Cartridge RAM
// 0xC000 - 0xCFFF : RAM Bank 0
// 0xD000 - 0xDFFF : RAM Bank 1-7 - switchable - Color only
// 0xE000 - 0xFDFF : Reserved - Echo RAM
// 0xFE00 - 0xFE9F : Object Attribute Memory
// 0xFEA0 - 0xFEFF : Reserved - Unusable
// 0xFF00 - 0xFF7F : I/O Registers
// 0xFF80 - 0xFFFE : Zero Page

type Bus struct {
	Cart *Cartridge
	RAM  *RAM
}

func NewBus(cart *Cartridge) *Bus {
	bus := Bus{Cart: cart, RAM: NewRAM()}

	return &bus
}

func (bus *Bus) Read(address uint16) byte {
	if address < 0x8000 {
		// ROM data
		value, err := bus.Cart.Read(address)
		if err != nil {
			fmt.Printf("Error reading from bus address %X: %s\n", address, err)
			return 0
		}
		return value
	} else if address < 0xE000 {
		//WRAM (Working RAM)
		return bus.RAM.WRAMRead(address)
	} else if address < 0xFF80 {
		//IO Registers
		return IORead(address)
	} else {
		fmt.Printf("Reading from bus address %X not implemented!\n", address)
		return 0
	}
}

func (bus *Bus) Read16(address uint16) uint16 {
	lo := bus.Read(address)
	hi := bus.Read(address + 1)

	return uint16(lo) | uint16(hi)<<8
}

func (bus *Bus) Write(address uint16, value byte) error {
	if address < 0x8000 {
		// ROM data
		err := bus.Cart.Write(address, value)
		if err != nil {
			return fmt.Errorf("Error writing to bus address %X: %s", address, err)
		}
		return nil
	} else if address < 0xE000 {
		//WRAM (Working RAM)
		bus.RAM.WRAMWrite(address, value)
		return nil
	} else if address < 0xFF80 {
		//IO Registers
		IOWrite(address, value)
		return nil
	} else {
		return fmt.Errorf("Writing to bus address %X not implemented!", address)
	}
}

func (bus *Bus) Write16(address uint16, value uint16) error {
	err := bus.Write(address+1, byte((value>>8)&0xFF))
	if err != nil {
		return err
	}
	err = bus.Write(address, byte(value&0xFF))
	if err != nil {
		return err
	}

	return nil
}
