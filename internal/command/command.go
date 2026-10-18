package command

import (
	"errors"
	"flag"
	"fmt"
	"io"
)

type commandType int

const (
	None commandType = iota
	Run
	Build
)

type Command struct {
	SourceFilePath string
	Command        commandType
	OutputFilePath string
	Verbose        bool
}

func (c *Command) Equals(other *Command) bool {
	if c.SourceFilePath != other.SourceFilePath ||
		c.Command != other.Command ||
		c.OutputFilePath != other.OutputFilePath ||
		c.Verbose != other.Verbose {
		return false
	}

	return true
}

func ParseArgs(args []string, stream io.Writer) (*Command, error) {
	arguments := Command{}

	if len(args) == 0 {
		PrintUsage(stream)
		return &arguments, nil
	}

	runFS := flag.NewFlagSet("run", flag.ContinueOnError)
	runFS.BoolVar(&arguments.Verbose, "verbose", false, "Verbose output")
	// TODO: verbose does nothing for now
	// TODO: -verbose after file argument does nothing

	buildFS := flag.NewFlagSet("build", flag.ContinueOnError)

	switch args[0] {
	case "run":
		runFS.Parse(args[1:])
		if len(runFS.Args()) < 1 {
			return &arguments, errors.New("Run requires an input file")
		}
		arguments.Command = Run
		arguments.SourceFilePath = runFS.Args()[0]
	case "build":
		buildFS.Parse(args[1:])
		if len(buildFS.Args()) < 1 {
			return &arguments, errors.New("Build requires an input file")
		}
		arguments.Command = Build
		arguments.SourceFilePath = buildFS.Args()[0]
	case "help", "-h", "-help", "--help":
		if len(args) > 1 {
			switch args[1] {
			case "run":
				runFS.Usage()
				arguments.Command = None
				return &arguments, nil
			case "build":
				buildFS.Usage()
				arguments.Command = None
				return &arguments, nil
			}
		}
		PrintUsage(stream)
		arguments.Command = None
	case "version", "-version", "--version":
		fmt.Println("bfc 0.0.1")
		arguments.Command = None
	default:
		return &arguments, errors.New("Unknown command")
	}
	return &arguments, nil
}

func PrintUsage(stream io.Writer) {
	fmt.Fprintf(stream, `Usage: bfc <command> [flags] <file>

Commands:
    run        interpret a Brainfuck program
    build      compile a Brainfuck program to an executable
    help       show help for a command
    version    show bfc version

Run 'bfc help <command>' for details`)
}
