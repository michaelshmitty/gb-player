package gb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInsert(t *testing.T) {
	cartridge, err := InsertCartridge("../roms/dmg-acid2.gb")

	assert := assert.New(t)

	assert.Nil(err)
	assert.Equal(cartridge.Filename, "../roms/dmg-acid2.gb")
	assert.Equal(cartridge.Title, "DMG-ACID2")
	assert.Equal(cartridge.Mapper, "ROM ONLY")
	assert.Equal(cartridge.Licensee, "None")
	assert.False(cartridge.SGBSupport, "SGB support should be false")
	assert.Equal(cartridge.ROMSize, 32)
	assert.Equal(cartridge.RAMSize, "0 - No RAM")
	assert.Equal(cartridge.Destination, "Japan (and possibly overseas)")
	assert.Equal(cartridge.Version, 0)
	assert.Equal(cartridge.Checksum, byte(0x9F))
}

func TestFailedChecksum(t *testing.T) {
	_, err := InsertCartridge("../roms/failed-checksum.gb")

	assert.EqualError(t, err, "ROM checksum failed: 9F does not equal 41")
}
