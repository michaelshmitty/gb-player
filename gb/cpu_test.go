package gb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func createCPU(data []byte) *CPU {
	cart := Cartridge{Filename: "test"}
	cart.Data = data
	cpu := NewCPU(NewBus(&cart))
	cpu.Regs.PC = 0

	return cpu
}

func TestSetFlag(t *testing.T) {
	cartridge, _ := InsertCartridge("../roms/dmg-acid2.gb")
	bus := NewBus(cartridge)
	cpu := NewCPU(bus)

	cpu.SetFlag(C)

	assert.True(t, cpu.IsFlagSet(C))
}

func TestClearFlag(t *testing.T) {
	cartridge, _ := InsertCartridge("../roms/dmg-acid2.gb")
	bus := NewBus(cartridge)
	cpu := NewCPU(bus)

	cpu.SetFlag(C)
	assert.True(t, cpu.IsFlagSet(C))

	cpu.ClearFlag(C)
	assert.False(t, cpu.IsFlagSet(C))
}

func TestToggleFlag(t *testing.T) {
	cartridge, _ := InsertCartridge("../roms/dmg-acid2.gb")
	bus := NewBus(cartridge)
	cpu := NewCPU(bus)

	cpu.ToggleFlag(C)
	assert.True(t, cpu.IsFlagSet(C))

	cpu.ToggleFlag(C)
	assert.False(t, cpu.IsFlagSet(C))
}

func TestInstruction00(t *testing.T) {
	cpu := createCPU([]byte{0x00, 0x00, 0x00})

	cpu.Step()

	assert.Equal(t, cpu.Regs.PC, uint16(0x01))
}

func TestInstruction3C(t *testing.T) {
	// Should increment A register
	cpu := createCPU([]byte{0x3C, 0x00, 0x00})

	assert.Equal(t, cpu.Regs.A, byte(0x0))
	cpu.Step()
	assert.Equal(t, cpu.Regs.A, byte(0x01))

	// Should clear Zero Flag
	cpu = createCPU([]byte{0x3C, 0x00, 0x00})

	cpu.SetFlag(Z)
	cpu.Step()
	assert.False(t, cpu.IsFlagSet(Z))

	// Should set Zero Flag
	cpu = createCPU([]byte{0x3C, 0x00, 0x00})
	cpu.Regs.A = 0xFF

	assert.False(t, cpu.IsFlagSet(Z))
	cpu.Step()
	assert.True(t, cpu.IsFlagSet(Z))

	// Should clear Subtract Flag
	cpu = createCPU([]byte{0x3C, 0x00, 0x00})

	cpu.SetFlag(N)
	cpu.Step()
	assert.False(t, cpu.IsFlagSet(N))

	// Should set Half Carry Flag if we overflow from bit 3
	cpu = createCPU([]byte{0x3C, 0x00, 0x00})
	cpu.Regs.A = 0x0F

	assert.False(t, cpu.IsFlagSet(H))
	cpu.Step()
	assert.True(t, cpu.IsFlagSet(H))

	// Should clear Half Carry Flag if we don't overflow from bit 3
	cpu = createCPU([]byte{0x3C, 0x00, 0x00})

	cpu.SetFlag(H)
	cpu.Step()
	assert.False(t, cpu.IsFlagSet(H))
}

func TestInstructionC3(t *testing.T) {
	cpu := createCPU([]byte{0xC3, 0x50, 0x01})

	cpu.Step()

	assert.Equal(t, cpu.Regs.PC, uint16(0x150))
}

func TestInstructionE9(t *testing.T) {
	cpu := createCPU([]byte{0xE9, 0x00, 0x00})
	cpu.Regs.H = 0x12
	cpu.Regs.L = 0x34

	cpu.Step()

	assert.Equal(t, cpu.Regs.PC, uint16(0x1234))
}
