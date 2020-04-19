package main

import (
	"bufio"
	"fmt"
	"glox/internal"
	"io/ioutil"
	"os"
)

func run(code []byte) (bool, error) {
	reporter := internal.StateErrorReporter{}
	frontend := internal.NewFrontend(code, &reporter)
	expr := frontend.Parse()

	if reporter.HadError {
		return reporter.HadError, nil
	}
	if expr != nil {
		printer := internal.Printer{}
		fmt.Println(printer.Print(expr))
	}
	return expr != nil, nil
}

func runFile(filePath string) error {
	if code, e := ioutil.ReadFile(filePath); e != nil {
		return e
	} else {
		if hadError, e := run(code); e != nil {
			return e
		} else if hadError {
			os.Exit(65)
		}
		return nil
	}
}

func runPrompt() error {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		if line, _, err := reader.ReadLine(); err != nil {
			return err
		} else if _, e := run(line); e != nil {
			// Run line and print error. The REPL shouldn't stop on a programming error,
			// unlike an IO error.
			fmt.Println(e)
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
