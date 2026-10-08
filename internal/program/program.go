package program

import "github.com/laurilyytikainen/go-bfc/internal/tokenparser"

type ProgramState struct {
	dataPointer        uint8
	instructionPointer uint
	program            []tokenparser.Token
}

func (p *ProgramState) Initialize(program []tokenparser.Token) {
	p.dataPointer = 0
	p.instructionPointer = 0
	p.program = program
}
