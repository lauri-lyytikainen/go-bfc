package tokenparser

import (
	"errors"
	"log"
	"os"
	"slices"
)

type Token int

const (
	tokenDataIncrement Token = iota
	tokenDataDecrement
	tokenByteIncrement
	tokenByteDecrement
	tokenOutputByte
	tokenInputByte
	tokenJumpIfZero
	tokenJumpToLoopStart
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
	case '<':
		o = tokenDataIncrement
	case '>':
		o = tokenDataDecrement
	case '+':
		o = tokenByteIncrement
	case '-':
		o = tokenByteDecrement
	case '[':
		o = tokenJumpIfZero
	case ']':
		o = tokenJumpToLoopStart
	case '.':
		o = tokenOutputByte
	case ',':
		o = tokenInputByte
	default:
		return 0, errors.New("Could not parse rune to token")
	}

	return o, nil
}
