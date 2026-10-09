package tokenparser

import (
	"errors"
	"log"
	"os"
	"slices"
)

type Token int

const (
	TokenDataIncrement Token = iota
	TokenDataDecrement
	TokenByteIncrement
	TokenByteDecrement
	TokenOutputByte
	TokenInputByte
	TokenJumpIfZero
	TokenJumpToLoopStart
)

func ParseSourceCode(fileName string) ([]Token, error) {
	bytes, err := os.ReadFile(fileName)
	if err != nil {
		return []Token{}, err
	}

	var program = []Token{}

	for i := range len(bytes) {
		currentRune := rune(bytes[i])
		if isProgramRune(currentRune) {
			token, err := parseToken(currentRune)
			if err != nil {
				log.Fatalf("%+v", err)
			}
			program = append(program, token)
		}
	}

	return program, nil
}

func isProgramRune(r rune) bool {
	allowedRunes := []rune{'<', '>', '+', '-', '[', ']', '.', ','}
	return slices.Contains(allowedRunes, r)
}

func parseToken(r rune) (Token, error) {
	var o Token = 0
	switch r {
	case '>':
		o = TokenDataIncrement
	case '<':
		o = TokenDataDecrement
	case '+':
		o = TokenByteIncrement
	case '-':
		o = TokenByteDecrement
	case '[':
		o = TokenJumpIfZero
	case ']':
		o = TokenJumpToLoopStart
	case '.':
		o = TokenOutputByte
	case ',':
		o = TokenInputByte
	default:
		return 0, errors.New("Could not parse rune to token")
	}

	return o, nil
}
