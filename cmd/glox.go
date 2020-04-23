package main

import (
	"bufio"
	"fmt"
	"glox/internal"
	"io/ioutil"
	"os"
)

type ErrorType int

const (
	HadNoError ErrorType = iota
	HadGeneralError
	HadRuntimeError
)

func run(code []byte) ErrorType {
	reporter := internal.StateErrorReporter{}
	frontend := internal.NewFrontend(code, &reporter)
	expr := frontend.Parse()
	interpreter := internal.NewInterpreter(&reporter)

	if reporter.HadError {
		return HadGeneralError
	}
	if expr != nil {
		interpreter.Interpret(expr)
		if reporter.HadRuntimeError {
			return HadRuntimeError
		}
	}
	return HadNoError
}

func runFile(filePath string) error {
	if code, e := ioutil.ReadFile(filePath); e != nil {
		return e
	} else {
		switch run(code) {
		case HadGeneralError:
			os.Exit(65)
		case HadRuntimeError:
			os.Exit(70)
		case HadNoError:
			return nil
		}
	}
	return nil
}

func runPrompt() error {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		if line, _, err := reader.ReadLine(); err != nil {
			return err
		} else {
			_ = run(line)
		}
	}
}

func main() {
	argv := os.Args[1:]

	if argc := len(argv); argc > 1 {
		fmt.Println("Usage: glox [script]")
		os.Exit(64)
	} else if argc == 1 {
		if e := runFile(argv[0]); e != nil {
			os.Exit(1)
		}
	} else {
		if e := runPrompt(); e != nil {
			fmt.Println(e)
			os.Exit(1)
		}
	}
}
