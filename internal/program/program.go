package program

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/laurilyytikainen/go-bfc/internal/tokenparser"
)

type ProgramState struct {
	dataPointer        uint
	instructionPointer uint
	program            []tokenparser.Token
	cells              map[uint]uint8
	stack              []uint
}

func (p *ProgramState) Initialize(program []tokenparser.Token) {
	p.dataPointer = 0
	p.instructionPointer = 0
	p.program = program
	p.cells = make(map[uint]uint8)
	p.stack = []uint{}
}

func (p *ProgramState) Execute() {
	for {
		if !p.processNextToken() {
			return
		}
	}
}

func (p *ProgramState) jumpToMatchingBrace(from uint) error {
	if p.program[p.instructionPointer] == tokenparser.TokenJumpIfZero {
		p.stack = append(p.stack, p.instructionPointer)
		for {
			if int(p.instructionPointer+1) < len(p.program) {
				p.instructionPointer++
				if p.program[p.instructionPointer] == tokenparser.TokenJumpIfZero {
					p.stack = append(p.stack, p.instructionPointer)
				} else if p.program[p.instructionPointer] == tokenparser.TokenJumpIfZero {
					if p.stack[len(p.stack)-1] == from {
						p.stack = p.stack[:len(p.stack)-1]
						break
					} else {
						p.stack = p.stack[:len(p.stack)-1]
					}
				}
			} else {
				return errors.New("Can not jump to matching brace")
			}
		}
	} else {
		return errors.New("Can not jump to matching brace from end of an program")
	}

	return nil
}

func (p *ProgramState) processNextToken() bool {
	if int(p.instructionPointer) > len(p.program) {
		return false
	}
	token := p.program[p.instructionPointer]
	switch token {
	case tokenparser.TokenDataIncrement:
		p.dataPointer = p.dataPointer + 1
	case tokenparser.TokenDataDecrement:
		if p.dataPointer > 0 {
			p.dataPointer = p.dataPointer - 1
		}
	case tokenparser.TokenByteIncrement:
		p.cells[p.dataPointer] = p.cells[p.dataPointer] + 1
	case tokenparser.TokenByteDecrement:
		p.cells[p.dataPointer] = p.cells[p.dataPointer] - 1
	case tokenparser.TokenOutputByte:
		fmt.Print(string(rune(p.cells[p.dataPointer])))
	case tokenparser.TokenInputByte:
		var inputString string = ""
		var value rune
		for {
			fmt.Println("Input: ")
			_, err := fmt.Scanln(&inputString)
			if err != nil {
				log.Fatalf("%+v", err)
				os.Exit(1)
			}
			if len(inputString) < 1 {
				fmt.Println("Please input at least one character")
				continue
			}
			value = rune(inputString[0])
			break
		}
		p.cells[p.dataPointer] = uint8(value)
	case tokenparser.TokenJumpIfZero:
		if p.cells[p.dataPointer] == 0 {
			p.jumpToMatchingBrace(p.instructionPointer)
			return true
		}
	case tokenparser.TokenJumpToLoopStart:
		return false
	default:
		return false
	}
	p.instructionPointer++
	return true
}
