package main

import (
	"log"
	"os"

	"github.com/laurilyytikainen/go-bfc/internal/program"
	"github.com/laurilyytikainen/go-bfc/internal/tokenparser"
)

func main() {
	code, err := tokenparser.ParseSourceCode("test.b")
	if err != nil {
		log.Fatalf("%+v", err)
	}

	program := program.ProgramState{}
	program.Initialize(code, os.Stdin, os.Stdin)

	program.Execute()
}
