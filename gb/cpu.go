package gb

import (
	"fmt"
	"os"
)

const CPUFrequency = 4194304

type CPUFlags byte

const (
	C CPUFlags = 1 << 4 // Carry Flag
	H CPUFlags = 1 << 5 // Half Carry Flag
	N CPUFlags = 1 << 6 // Subtract Flag
	Z CPUFlags = 1 << 7 // Zero Flag
)

type Registers struct {
	A  byte
	F  CPUFlags
	B  byte
	C  byte
	D  byte
	E  byte
	H  byte
	L  byte
	PC uint16
	SP uint16
}

type CPU struct {
	Bus      *Bus
	Regs     Registers
	Halted   bool
	Stepping bool
}

func NewCPU(bus *Bus) *CPU {
	cpu := CPU{}
	cpu.Bus = bus
	cpu.Regs = Registers{}
	cpu.Stepping = true

	return &cpu
}

func (cpu *CPU) Step() {
	if !cpu.Halted {
		opcode := cpu.Bus.Read(cpu.Regs.PC)
		cpu.Regs.PC++

		fmt.Printf("%04X: (%02X %02X %02X) A: %02X B: %02X C: %02X\n", cpu.Regs.PC,
			opcode, cpu.Bus.Read(cpu.Regs.PC), cpu.Bus.Read(cpu.Regs.PC+1), cpu.Regs.A, cpu.Regs.B, cpu.Regs.C)

		switch opcode {

		case 0x00:
			// NOP
		case 0x3C:
			// INC A
			cpu.Regs.A++

			// Set appropriate flags
			if cpu.Regs.A == 0 {
				cpu.SetFlag(Z)
			} else {
				cpu.ClearFlag(Z)
			}

			cpu.ClearFlag(N)

			if (cpu.Regs.A & 0x0F) == 0 {
				cpu.SetFlag(H)
			} else {
				cpu.ClearFlag(H)
			}
		case 0xC3:
			// JP a16
			lo := cpu.Bus.Read(cpu.Regs.PC)
			//     emu_cycles(1);
			hi := cpu.Bus.Read(cpu.Regs.PC + 1)
			// emu_cycles(1);
			cpu.Regs.PC = uint16(hi)<<8 | uint16(lo)

		case 0xE9:
			// JP HL
			val := uint16(cpu.Regs.H)<<8 | uint16(cpu.Regs.L)
			cpu.Regs.PC = val

		default:
			fmt.Printf("\nINVALID INSTRUCTION! Unknown opcode: %02X\n", opcode)
			os.Exit(1)
		}
	}
}

func (cpu *CPU) SetFlag(flag CPUFlags) {
	cpu.Regs.F |= flag
}

func (cpu *CPU) ClearFlag(flag CPUFlags) {
	cpu.Regs.F &^= flag
}

func (cpu *CPU) ToggleFlag(flag CPUFlags) {
	cpu.Regs.F ^= flag
}

func (cpu *CPU) IsFlagSet(flag CPUFlags) bool {
	return cpu.Regs.F&flag != 0
}

func (cpu *CPU) StackPush(data byte) {
	cpu.Regs.SP--
	cpu.Bus.Write(cpu.Regs.SP, data)
}

func (cpu *CPU) StackPush16(data uint16) {
	cpu.StackPush(byte((data >> 8) & 0xFF))
	cpu.StackPush(byte(data & 0xFF))
}

func (cpu *CPU) StackPop() byte {
	val := cpu.Bus.Read(cpu.Regs.SP)
	cpu.Regs.SP++

	return val
}

func (cpu *CPU) StackPop16() uint16 {
	lo := cpu.StackPop()
	hi := cpu.StackPop()

	return uint16(hi)<<8 | uint16(lo)
}
