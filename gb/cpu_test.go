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

	assert.Equal(t, uint16(0x01), cpu.Regs.PC)
}

func TestInstruction3C(t *testing.T) {
	// Should increment A register
	cpu := createCPU([]byte{0x3C, 0x00, 0x00})

	assert.Equal(t, byte(0x0), cpu.Regs.A)
	cpu.Step()
	assert.Equal(t, byte(0x01), cpu.Regs.A)

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

	assert.Equal(t, uint16(0x150), cpu.Regs.PC)
}

func TestInstructionE9(t *testing.T) {
	cpu := createCPU([]byte{0xE9, 0x00, 0x00})
	cpu.Regs.H = 0x12
	cpu.Regs.L = 0x34

	cpu.Step()

	assert.Equal(t, uint16(0x1234), cpu.Regs.PC)
}

func TestInstructionF3(t *testing.T) {
	cpu := createCPU([]byte{0xF3, 0x00, 0x00})

	cpu.InterruptMasterEnable = true
	cpu.Step()

	assert.False(t, cpu.InterruptMasterEnable)
}

func TestInstruction31(t *testing.T) {
	cpu := createCPU([]byte{0x31, 0xFF, 0xCF})

	cpu.Step()

	// Should load SP with n16 value
	assert.Equal(t, uint16(0xCFFF), cpu.Regs.SP)

	// Should step over the 16 bit value onto the next instruction, i.e. increase the program counter with 3.
	assert.Equal(t, uint16(0x0003), cpu.Regs.PC)
}

func TestInstructionCD(t *testing.T) {
	cpu := createCPU([]byte{0xCD, 0x4C, 0x48, 0x41})

	cpu.Step()

	// Should push the address of the next instruction onto the stack
	// In this case 0x41 which is at address 0x03 (i.e. the 4th instruction in the row above creating the CPU instance)
	assert.Equal(t, uint16(0x03), cpu.Bus.Read16(cpu.Regs.SP))

	// Should jump to n16 value
	assert.Equal(t, uint16(0x484C), cpu.Regs.PC)
}

func TestInstruction21(t *testing.T) {
	cpu := createCPU([]byte{0x21, 0x34, 0x12})

	cpu.Step()

	// Should load the 16-bit immediate value into the combined HL register
	assert.Equal(t, byte(0x12), cpu.Regs.H)
	assert.Equal(t, byte(0x34), cpu.Regs.L)

	// Should increase the stack pointer
	assert.Equal(t, uint16(0x3), cpu.Regs.PC)
}
