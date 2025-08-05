package gb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInsertCartridge(t *testing.T) {
	cartridge := Insert("../rom.gb")
	assert := assert.New(t)

	assert.Equal(cartridge.Title, "SEIKEN DENSETSU")
	assert.Equal(cartridge.Mapper, "MBC2+BATTERY")
	assert.Equal(cartridge.Licensee, "Square")
	assert.False(cartridge.SGBSupport, "SGB support should be false")
	assert.Equal(cartridge.ROMSize, 256)
	assert.Equal(cartridge.RAMSize, "0 - No RAM")
	assert.Equal(cartridge.Destination, "Overseas only")
	assert.Equal(cartridge.Version, 0)
}
