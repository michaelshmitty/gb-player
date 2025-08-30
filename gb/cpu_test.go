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

// func TestInstructionCB7E(t *testing.T) {
// 	cpu := createCPU([]byte{0xCB, 0x7E, 0x00, 0x00})

// 	cpu.Regs.H = 0xFF
// 	cpu.Regs.L = 0x40
// 	err := cpu.Bus.Write(0xFF40, 0xFF)
// 	assert.Equal(t, err, nil)

// 	// FIXME(m): This needs bus access to IO, in particular the LCD
// 	// at 0xFF40.

// 	// assert.Equal(t, 0xFF, cpu.Bus.Read(0xFF40))

// 	// cpu.Step()

// 	// // Should set the zero flag
// 	// assert.True(t, cpu.IsFlagSet(Z))
// }

func TestInstruction47(t *testing.T) {
	cpu := createCPU([]byte{0x47, 0x00, 0x00})

	cpu.Regs.A = 0xDE

	cpu.Step()

	// Should load value in register A into register B
	assert.Equal(t, byte(0xDE), cpu.Regs.B)
}

func TestInstruction11(t *testing.T) {
	cpu := createCPU([]byte{0x11, 0x34, 0x12})

	cpu.Step()

	// Should load the 16-bit immediate value into the combined DE register
	assert.Equal(t, byte(0x12), cpu.Regs.D)
	assert.Equal(t, byte(0x34), cpu.Regs.E)

	// Should increase the stack pointer
	assert.Equal(t, uint16(0x3), cpu.Regs.PC)
}

func TestInstruction0E(t *testing.T) {
	cpu := createCPU([]byte{0x0E, 0xDE, 0x00})

	cpu.Step()

	// Should load the 8-bit immediate value into register C
	assert.Equal(t, byte(0xDE), cpu.Regs.C)

	// Should increase the stack pointer
	assert.Equal(t, uint16(0x2), cpu.Regs.PC)
}

func TestInstruction2A(t *testing.T) {
	cpu := createCPU([]byte{0x2A, 0x00, 0x00})
	cpu.Regs.H = 0xD1
	cpu.Regs.L = 0x23
	cpu.Bus.Write(0xD123, 0xDE)

	cpu.Step()

	// Should load the byte pointed to by HL into register A
	assert.Equal(t, byte(0xDE), cpu.Regs.A)

	// Should increment HL
	hl := uint16(cpu.Regs.H)<<8 | uint16(cpu.Regs.L)
	assert.Equal(t, uint16(0xD124), hl)

	// Should increase the stack pointer
	assert.Equal(t, uint16(0x1), cpu.Regs.PC)
}

func TestInstruction12(t *testing.T) {
	cpu := createCPU([]byte{0x12, 0x00, 0x00})
	cpu.Regs.D = 0xD1
	cpu.Regs.E = 0x23
	cpu.Regs.A = 0xDE

	cpu.Step()

	// Should copy the value in register A into the byte pointed to by DE
	address := uint16(cpu.Regs.D)<<8 | uint16(cpu.Regs.E)
	assert.Equal(t, byte(0xDE), cpu.Bus.Read(address))

	// Should increase the stack pointer
	assert.Equal(t, uint16(0x1), cpu.Regs.PC)
}

func TestInstruction1C(t *testing.T) {
	// Should increment E register
	cpu := createCPU([]byte{0x1C, 0x00, 0x00})

	cpu.Step()
	assert.Equal(t, byte(0x01), cpu.Regs.E)

	// Should clear Zero Flag
	cpu = createCPU([]byte{0x1C, 0x00, 0x00})

	cpu.SetFlag(Z)
	cpu.Step()
	assert.False(t, cpu.IsFlagSet(Z))

	// Should set Zero Flag
	cpu = createCPU([]byte{0x1C, 0x00, 0x00})
	cpu.Regs.E = 0xFF

	cpu.Step()
	assert.True(t, cpu.IsFlagSet(Z))

	// Should clear Subtract Flag
	cpu = createCPU([]byte{0x1C, 0x00, 0x00})

	cpu.SetFlag(N)
	cpu.Step()
	assert.False(t, cpu.IsFlagSet(N))

	// Should set Half Carry Flag if we overflow from bit 3
	cpu = createCPU([]byte{0x1C, 0x00, 0x00})
	cpu.Regs.E = 0x0F

	cpu.Step()
	assert.True(t, cpu.IsFlagSet(H))

	// Should clear Half Carry Flag if we don't overflow from bit 3
	cpu = createCPU([]byte{0x1C, 0x00, 0x00})

	cpu.SetFlag(H)
	cpu.Step()
	assert.False(t, cpu.IsFlagSet(H))
}

func TestInstruction20(t *testing.T) {
	// Should not jump if Z is set
	cpu := createCPU([]byte{0x20, 0xFB, 0x00})
	cpu.SetFlag(Z)

	cpu.Step()

	assert.Equal(t, uint16(0x02), cpu.Regs.PC)

	// Should jump to positive offset if Z is not set
	cpu = createCPU([]byte{0x20, 0x0A, 0x00})

	cpu.Step()

	assert.Equal(t, uint16(0x0B), cpu.Regs.PC)

	// Should jump to negative offset if Z is not set
	cpu = createCPU([]byte{0x20, 0xFB, 0x00})

	cpu.Step()

	assert.Equal(t, uint16(0xFFFC), cpu.Regs.PC)
}

func TestInstruction10(t *testing.T) {
	cpu := createCPU([]byte{0x10, 0x2A, 0x12})

	cpu.Step()

	assert.True(t, cpu.Halted)
}
