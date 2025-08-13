package gb

import (
	"fmt"
)

var instructions = map[byte]string{
	0x00: "NOP",
	0x01: "FOO",
}

type Instruction struct {
}

func InstructionByOpcode(opcode byte) (string, error) {
	instruction, ok := instructions[opcode]
	if !ok {
		return "", fmt.Errorf("Unknown opcode: %02X", opcode)
	}

	return instruction, nil
}
