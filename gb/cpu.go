package gb

import (
	"fmt"
)

const CPUFrequency = 4194304

type Registers struct {
	A  byte
	F  byte
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
	Bus                    *Bus
	Regs                   Registers
	FetchedData            uint16
	MemoryDestination      uint16
	DestinationIsMemory    bool
	CurrentOpcode          byte
	Halted                 bool
	Stepping               bool
	InterruptMasterEnabled bool
	CurrentInstruction     string
}

func NewCPU(bus *Bus) *CPU {
	cpu := CPU{}
	cpu.Bus = bus
	cpu.Regs = Registers{PC: 0x100}
	cpu.Stepping = true

	return &cpu
}

func (cpu *CPU) Step() error {
	if !cpu.Halted {
		err := cpu.fetchInstruction()
		if err != nil {
			return fmt.Errorf("Error fetching instruction: %s", err)
		}
		cpu.fetchData()
		cpu.execute()
	}

	return nil
}

func (cpu *CPU) fetchInstruction() error {
	opcode, err := cpu.Bus.Read(cpu.Regs.PC)
	if err != nil {
		return fmt.Errorf("Error fetching instruction at address %X", cpu.Regs.PC)
	}
	cpu.CurrentOpcode = opcode
	cpu.CurrentInstruction, err = InstructionByOpcode(cpu.CurrentOpcode)
	if err != nil {
		return fmt.Errorf("Error translating opcode %02X: %s at PC: %04X", cpu.CurrentOpcode, err, cpu.Regs.PC)
	}

	fmt.Printf("Executing instruction: %02X	PC: %04X\n", cpu.CurrentOpcode, cpu.Regs.PC)

	cpu.Regs.PC++

	return nil
}

func (cpu *CPU) fetchData() {

}

func (cpu *CPU) execute() {
	fmt.Println("Not executing yet...")
}
