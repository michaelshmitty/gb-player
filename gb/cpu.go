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
	Bus                   *Bus
	Regs                  Registers
	Halted                bool
	Stepping              bool
	InterruptMasterEnable bool
	InterruptFlags        byte
}

func NewCPU(bus *Bus) *CPU {
	cpu := CPU{}
	cpu.Bus = bus
	// NOTE(m): PC is usually set to 0x100 by the boot rom
	// TODO(m): SP is usually set programmatically by the cartridge code.
	// Remove this hardcoded value later!
	cpu.Regs = Registers{SP: 0xDFFF}
	cpu.Stepping = true

	return &cpu
}

func (cpu *CPU) Step() {
	if !cpu.Halted {
		opcode := cpu.Bus.Read(cpu.Regs.PC)

		fmt.Printf("%04X: (%02X %02X %02X) A: %02X B: %02X C: %02X D: %02X E: %02X H: %02X L: %02X\n", cpu.Regs.PC,
			opcode, cpu.Bus.Read(cpu.Regs.PC+1), cpu.Bus.Read(cpu.Regs.PC+2), cpu.Regs.A, cpu.Regs.B, cpu.Regs.C,
			cpu.Regs.D, cpu.Regs.E, cpu.Regs.H, cpu.Regs.L)

		cpu.Regs.PC++

		switch opcode {

		case 0x00:
			// NOP

		case 0x0D:
			// DEC C
			cpu.Regs.C--

			// Set appropriate flags
			if cpu.Regs.C == 0 {
				cpu.SetFlag(Z)
			} else {
				cpu.ClearFlag(Z)
			}

			cpu.SetFlag(N)

			if (cpu.Regs.C & 0x0F) == 0x0F {
				cpu.SetFlag(H)
			} else {
				cpu.ClearFlag(H)
			}

		case 0x0E:
			// LD C, n8
			val := cpu.Bus.Read(cpu.Regs.PC)
			//     emu_cycles(1);
			cpu.Regs.C = val
			cpu.Regs.PC++

		case 0x10:
			// STOP n8
			cpu.Halted = true

		case 0x11:
			// LD DE, n16
			cpu.Regs.E = cpu.Bus.Read(cpu.Regs.PC)
			//     emu_cycles(1);
			cpu.Regs.D = cpu.Bus.Read(cpu.Regs.PC + 1)
			// emu_cycles(1);
			cpu.Regs.PC += 2

		case 0x12:
			// LD [DE], A
			// Get 16-bit address from DE
			address := uint16(cpu.Regs.D)<<8 | uint16(cpu.Regs.E)
			cpu.Bus.Write(address, cpu.Regs.A)

		case 0x14:
			// INC D
			cpu.Regs.D++

			// Set appropriate flags
			if cpu.Regs.D == 0 {
				cpu.SetFlag(Z)
			} else {
				cpu.ClearFlag(Z)
			}

			cpu.ClearFlag(N)

			if (cpu.Regs.D & 0x0F) == 0 {
				cpu.SetFlag(H)
			} else {
				cpu.ClearFlag(H)
			}

		case 0x18:
			// JR e8
			// Jump relative to 8-bit signed offset
			// emu_cycles(3);
			offset := int8(cpu.Bus.Read(cpu.Regs.PC))
			cpu.Regs.PC++
			cpu.Regs.PC = uint16(int(cpu.Regs.PC) + int(offset))

		case 0x1C:
			// INC E
			cpu.Regs.E++

			// Set appropriate flags
			if cpu.Regs.E == 0 {
				cpu.SetFlag(Z)
			} else {
				cpu.ClearFlag(Z)
			}

			cpu.ClearFlag(N)

			if (cpu.Regs.E & 0x0F) == 0 {
				cpu.SetFlag(H)
			} else {
				cpu.ClearFlag(H)
			}

		case 0x20:
			// JR NZ, e8
			if cpu.IsFlagSet(Z) {
				// Z is set, don't execute
				// emu_cycles(2);
				cpu.Regs.PC++
			} else {
				// Z is not set, execute
				// Jump relative to 8-bit signed offset
				// emu_cycles(3);
				offset := int8(cpu.Bus.Read(cpu.Regs.PC))
				cpu.Regs.PC++
				cpu.Regs.PC = uint16(int(cpu.Regs.PC) + int(offset))
			}

		case 0x21:
			// LD HL, n16
			cpu.Regs.L = cpu.Bus.Read(cpu.Regs.PC)
			//     emu_cycles(1);
			cpu.Regs.H = cpu.Bus.Read(cpu.Regs.PC + 1)
			// emu_cycles(1);
			cpu.Regs.PC += 2

		case 0x2A:
			// LD A, [HL+] or LD A, [HLI]
			// Get 16-bit address from HL
			address := uint16(cpu.Regs.H)<<8 | uint16(cpu.Regs.L)
			// Read byte at address and assign value to A register
			cpu.Regs.A = cpu.Bus.Read(address)
			// emu_cycles(1);
			// Increment HL
			address++
			cpu.Regs.H = byte(address >> 8)
			cpu.Regs.L = byte(address)
			// emu_cycles(1);

		case 0x31:
			// LD SP, n16
			lo := cpu.Bus.Read(cpu.Regs.PC)
			//     emu_cycles(1);
			hi := cpu.Bus.Read(cpu.Regs.PC + 1)
			// emu_cycles(1);
			cpu.Regs.SP = uint16(lo) | uint16(hi)<<8
			cpu.Regs.PC += 2

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

		case 0x3E:
			// LD A, n8
			val := cpu.Bus.Read(cpu.Regs.PC)
			//     emu_cycles(1);
			cpu.Regs.A = val
			cpu.Regs.PC++

		case 0x47:
			// LD B, A
			cpu.Regs.B = cpu.Regs.A

		case 0x78:
			// LD A, B
			cpu.Regs.A = cpu.Regs.B

		case 0x7C:
			// LD A, H
			cpu.Regs.A = cpu.Regs.H

		case 0x7D:
			// LD A, L
			cpu.Regs.A = cpu.Regs.L

		case 0xCB:
			// Prefix byte instructions
			cbOpcode := cpu.Bus.Read(cpu.Regs.PC)

			fmt.Printf("%04X: (%02X %02X %02X) A: %02X B: %02X C: %02X\n", cpu.Regs.PC,
				cbOpcode, cpu.Bus.Read(cpu.Regs.PC+1), cpu.Bus.Read(cpu.Regs.PC+2), cpu.Regs.A, cpu.Regs.B, cpu.Regs.C)

			cpu.Regs.PC++

			switch cbOpcode {
			// case 0x7E:
			// 	// BIT 7, [HL]
			// 	// Read byte pointed to by address HL
			// 	address := uint16(cpu.Regs.H)<<8 | uint16(cpu.Regs.L)
			// 	val := cpu.Bus.Read(address)

			// 	// Check if bit 7 is set
			// 	if (val & 0x80) == 0 {
			// 		// Set zero flag if bit is not set
			// 		cpu.SetFlag(Z)
			// 	}
			// 	cpu.ClearFlag(N)
			// 	cpu.SetFlag(H)

			default:
				fmt.Printf("\nINVALID INSTRUCTION! Unknown CB opcode: %02X\n", cbOpcode)
				os.Exit(1)
			}

		case 0xC3:
			// JP a16
			lo := cpu.Bus.Read(cpu.Regs.PC)
			//     emu_cycles(1);
			hi := cpu.Bus.Read(cpu.Regs.PC + 1)
			// emu_cycles(1);
			cpu.Regs.PC = uint16(lo) | uint16(hi)<<8

		case 0xC9:
			// RET
			// emu_cycles(4);
			cpu.Regs.PC = cpu.StackPop16()

		case 0xCD:
			// CALL a16
			cpu.StackPush16(cpu.Regs.PC + 2)

			lo := cpu.Bus.Read(cpu.Regs.PC)
			//     emu_cycles(1);
			hi := cpu.Bus.Read(cpu.Regs.PC + 1)
			// emu_cycles(1);
			cpu.Regs.PC = uint16(lo) | uint16(hi)<<8

		case 0xE0:
			// LDH [a8], A
			offset := cpu.Bus.Read(cpu.Regs.PC)
			address := 0xFF00 | uint16(offset)
			cpu.Bus.Write(address, cpu.Regs.A)
			cpu.Regs.PC++

		case 0xE5:
			// PUSH HL
			// emu_cycles(4);
			cpu.StackPush(cpu.Regs.H)
			cpu.StackPush(cpu.Regs.L)

		case 0xE9:
			// JP HL
			val := uint16(cpu.Regs.H)<<8 | uint16(cpu.Regs.L)
			cpu.Regs.PC = val

		case 0xEA:
			// LD [a16], A
			lo := cpu.Bus.Read(cpu.Regs.PC)
			cpu.Regs.PC++
			hi := cpu.Bus.Read(cpu.Regs.PC)
			cpu.Regs.PC++

			address := uint16(lo) | uint16(hi)<<8
			cpu.Bus.Write(address, cpu.Regs.A)

		case 0xF3:
			// DI
			cpu.InterruptMasterEnable = false

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

func (cpu *CPU) GetInterruptFlags() byte {
	return cpu.InterruptFlags
}

func (cpu *CPU) SetInterruptFlags(value byte) {
	cpu.InterruptFlags = value
}
