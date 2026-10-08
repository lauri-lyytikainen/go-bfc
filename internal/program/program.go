package program

import (
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
}

func (p *ProgramState) Initialize(program []tokenparser.Token) {
	p.dataPointer = 0
	p.instructionPointer = 0
	p.program = program
	p.cells = make(map[uint]uint8)
}

func (p *ProgramState) Execute() {
	for {
		if !p.processNextToken() {
			return
		}
	}
}

func (p *ProgramState) processNextToken() bool {
	if int(p.instructionPointer) > len(p.program) {
		return false
	}
	token := p.program[p.instructionPointer]
	switch token {
	case tokenparser.TokenDataIncrement:
		log.Println("Token: >")
		p.dataPointer++
	case tokenparser.TokenDataDecrement:
		log.Println("Token: <")
		p.dataPointer--
	case tokenparser.TokenByteIncrement:
		log.Println("Token: +")
		p.cells[p.dataPointer] = p.cells[p.dataPointer] + 1
	case tokenparser.TokenByteDecrement:
		log.Println("Token: -")
		p.cells[p.dataPointer] = p.cells[p.dataPointer] - 1
	case tokenparser.TokenOutputByte:
		log.Println("Token: .")
		fmt.Print(string(rune(p.cells[p.dataPointer])))
	case tokenparser.TokenInputByte:
		log.Println("Token: ,")
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
		log.Println("Token: [")
		return false
	case tokenparser.TokenJumpToLoopStart:
		log.Println("Token: ]")
		return false
	default:
		return false
	}
	p.instructionPointer++
	return true
}
