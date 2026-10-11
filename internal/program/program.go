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
	jumpTable          map[uint]uint
}

func (p *ProgramState) Initialize(program []tokenparser.Token) {
	p.dataPointer = 0
	p.instructionPointer = 0
	p.program = program
	p.cells = make(map[uint]uint8)
	err := p.createJumpTable()
	if err != nil {
		log.Fatalf("%+v", err)
		os.Exit(1)
	}
}

func (p *ProgramState) createJumpTable() error {
	p.jumpTable = make(map[uint]uint)
	stack := []uint{}
	for index, token := range p.program {
		if token == tokenparser.TokenJumpIfZero {
			stack = append(stack, uint(index))
			continue
		} else if token == tokenparser.TokenJumpToLoopStart {
			if len(stack) == 0 {
				return errors.New("Unmatched ']'")
			} else {
				p.jumpTable[stack[len(stack)-1]] = uint(index)
				p.jumpTable[uint(index)] = stack[len(stack)-1]
				stack = stack[:len(stack)-1]
			}
		}
	}
	if len(stack) > 0 {
		return errors.New("Unmatched '['")
	}

	return nil
}

func (p *ProgramState) Execute() {
	for {
		if !p.processNextToken() {
			return
		}
	}
}

func (p *ProgramState) processNextToken() bool {
	if int(p.instructionPointer) >= len(p.program) {
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
			_, err := fmt.Scanln(&inputString)
			// TODO: Read one byte at a time using bufio reader and readbyte
			if err != nil {
				log.Fatalf("%+v", err)
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
			target, ok := p.jumpTable[p.instructionPointer]
			if ok {
				p.instructionPointer = target
			} else {
				log.Fatalf("Failed to jump to mathing brace, not found in jump list")
			}
		}
	case tokenparser.TokenJumpToLoopStart:
		if p.cells[p.dataPointer] != 0 {
			target, ok := p.jumpTable[p.instructionPointer]
			if ok {
				p.instructionPointer = target
			} else {
				log.Fatalf("Failed to jump to mathing brace, not found in jump list")
			}
		}
	default:
		return false
	}
	p.instructionPointer++
	return true
}
