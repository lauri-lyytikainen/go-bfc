package main

import (
	"fmt"
	"log"
	"os"

	"github.com/laurilyytikainen/go-bfc/internal/command"
	"github.com/laurilyytikainen/go-bfc/internal/program"
	"github.com/laurilyytikainen/go-bfc/internal/tokenparser"
)

func main() {

	osArgs := os.Args[1:]
	args, err := command.ParseArgs(osArgs, os.Stdout)
	if err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(1)
	}

	switch args.Command {
	case command.Run:
		code, err := tokenparser.ParseSourceCode(args.SourceFilePath)
		if err != nil {
			log.Fatalf("%+v", err)
		}

		program := program.ProgramState{}
		program.Initialize(code, os.Stdin, os.Stdout)

		program.Execute()
	case command.Build:
		fmt.Println("Building is not implemented yet")
	default:
	}
}


