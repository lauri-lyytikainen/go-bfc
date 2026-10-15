package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/laurilyytikainen/go-bfc/internal/program"
	"github.com/laurilyytikainen/go-bfc/internal/tokenparser"
)

type commandType int

const (
	Run commandType = iota
	Build
)

type Command struct {
	sourceFilePath string
	command        commandType
	outputFilePath string
	verbose        bool
}

func main() {

	osArgs := os.Args[1:]
	args, err := parseArgs(osArgs)
	if err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(1)
	}

	switch args.command {
	case Run:
		code, err := tokenparser.ParseSourceCode(args.sourceFilePath)
		if err != nil {
			log.Fatalf("%+v", err)
		}

		program := program.ProgramState{}
		program.Initialize(code, os.Stdin, os.Stdout)

		program.Execute()
	case Build:
		fmt.Println("Building is not implemented yet")
	}
}

func parseArgs(args []string) (*Command, error) {
	arguments := Command{}

	if len(args) == 0 {
		printUsage()
		// TODO: improve the structure of this file
		return &arguments, errors.New("usage error")
	}

	runFS := flag.NewFlagSet("run", flag.ContinueOnError)
	runFS.BoolVar(&arguments.verbose, "verbose", false, "Verbose output")

	buildFS := flag.NewFlagSet("build", flag.ContinueOnError)

	switch args[0] {
	case "run":
		runFS.Parse(args[1:])
		if len(runFS.Args()) < 1 {
			return &arguments, errors.New("Run requires an input file")
		}
		arguments.command = Run
		arguments.sourceFilePath = runFS.Args()[0]
	case "build":
		buildFS.Parse(args[1:])
		if len(buildFS.Args()) < 1 {
			return &arguments, errors.New("Build requires an input file")
		}
		arguments.command = Build
		arguments.sourceFilePath = buildFS.Args()[0]
	case "help", "-h", "-help", "--help":
		if len(args) > 1 {
			switch args[1] {
			case "run":
				runFS.Usage()
				os.Exit(0)
			case "build":
				buildFS.Usage()
				os.Exit(0)
			}
		}
		printUsage()
		os.Exit(0)
	case "version", "-version", "--version":
		fmt.Println("bfc 0.0.1")
		os.Exit(0)
	default:
		return &arguments, errors.New("Unknown command")
	}
	return &arguments, nil
}

func printUsage() {
	fmt.Fprintf(os.Stderr, `Usage: bfc <command> [flags] <file>

Commands:
    run        interpret a Brainfuck program
    build      compile a Brainfuck program to an executable
    help       show help for a command
    version    show bfc version

Run 'bfc help <command>' for details`)
}
