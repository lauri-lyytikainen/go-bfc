package program

import (
	"bufio"
	"errors"
	"io"
	"log"

	"github.com/laurilyytikainen/go-bfc/internal/tokenparser"
)

type ProgramState struct {
	dataPointer        uint
	instructionPointer uint
	program            []tokenparser.Token
	cells              map[uint]uint8
	jumpTable          map[uint]uint
	reader             *bufio.Reader
	writer             *bufio.Writer
}

func (p *ProgramState) Initialize(
	program []tokenparser.Token,
	ioReader io.Reader,
	ioWriter io.Writer,
) {
	p.dataPointer = 0
	p.instructionPointer = 0
	p.program = program
	p.cells = make(map[uint]uint8)
	err := p.createJumpTable()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	p.reader = bufio.NewReader(ioReader)
	p.writer = bufio.NewWriter(ioWriter)
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
		p.writer.WriteByte(p.cells[p.dataPointer])
	case tokenparser.TokenInputByte:
		err := p.writer.Flush()
		if err != nil {
			log.Fatalf("Failed to flush output")
		}
		byte, err := p.reader.ReadByte()
		if err != nil {
			if err != io.EOF {
				log.Fatalf("Failed to read byte from input")
			}
			p.cells[p.dataPointer] = 0
		} else {
			p.cells[p.dataPointer] = byte
		}
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
	err := p.writer.Flush()
	if err != nil {
		log.Fatalf("Failed to flush output")
	}
	return true
}
