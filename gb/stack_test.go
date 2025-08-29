package gb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStackPush(t *testing.T) {
	cpu := createCPU([]byte{0x00, 0x00, 0x00})

	cpu.StackPush(0xDE)

	// Should decrement the stack pointer
	assert.Equal(t, cpu.Regs.SP, uint16(0xDFFE))

	// Should write the value to the stack
	assert.Equal(t, cpu.Bus.Read(cpu.Regs.SP), byte(0xDE))
}

func TestStackPush16(t *testing.T) {
	cpu := createCPU([]byte{0x00, 0x00, 0x00})

	cpu.StackPush16(0xDEAD)

	// Should decrement the stack pointer twice
	assert.Equal(t, cpu.Regs.SP, uint16(0xDFFD))

	// Should write the value to the stack
	assert.Equal(t, cpu.Bus.Read16(cpu.Regs.SP), uint16(0xDEAD))
}

func TestStackPop(t *testing.T) {
	cpu := createCPU([]byte{0x00, 0x00, 0x00})

	cpu.StackPush(0xDE)

	val := cpu.StackPop()

	// Should increment the stack pointer
	assert.Equal(t, cpu.Regs.SP, uint16(0xDFFF))

	// Should return the byte value
	assert.Equal(t, val, byte(0xDE))
}

func TestStackPop16(t *testing.T) {
	cpu := createCPU([]byte{0x00, 0x00, 0x00})

	cpu.StackPush16(0xBEEF)

	val := cpu.StackPop16()

	// Should increment the stack pointer
	assert.Equal(t, cpu.Regs.SP, uint16(0xDFFF))

	// Should return the 16 bit value
	assert.Equal(t, val, uint16(0xBEEF))
}
