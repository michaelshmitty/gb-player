package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInsertCartridge(t *testing.T) {
	cartridge := Insert("./roms/dmg-acid2.gb")
	assert := assert.New(t)

	assert.Equal(cartridge.Title, "DMG-ACID2")
	assert.Equal(cartridge.Mapper, "ROM ONLY")
	assert.Equal(cartridge.Licensee, "None")
	assert.False(cartridge.SGBSupport, "SGB support should be false")
	assert.Equal(cartridge.ROMSize, 32)
	assert.Equal(cartridge.RAMSize, "0 - No RAM")
	assert.Equal(cartridge.Destination, "Japan (and possibly overseas)")
	assert.Equal(cartridge.Version, 0)
}
