package gb

type Cpu struct {
	A     uint8
	Flags uint8
	BC    uint16
	DE    uint16
	HL    uint16
	SP    uint16
	PC    uint16
}

func Reset() Cpu {
	cpu := Cpu{}

	return cpu
}

func Tick(cpu *Cpu) {
	cpu.PC++
}
