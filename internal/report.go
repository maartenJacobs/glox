package internal

import (
	"fmt"
	"os"
)

// ErrorReporter provides a simple error reporting service that can be shared between
// different parts of the compiler.
type ErrorReporter interface {
	Error(line int, message string)
	Report(line int, where string, message string)
	RuntimeError(e RuntimeError)
}

// StateErrorReporter is an implementation of ErrorReporter that tracks whether an
// error was reported and prints errors to standard error.
type StateErrorReporter struct {
	HadError        bool // Whether an error has been reported.
	HadRuntimeError bool // Whether a runtime error has been thrown.
}

func (reporter *StateErrorReporter) Error(line int, message string) {
	reporter.Report(line, "", message)
}

func (reporter *StateErrorReporter) Report(line int, where string, message string) {
	_, err := fmt.Fprintf(os.Stderr, "[line %d] Error%s: %s\n", line, where, message)
	if err != nil { // Not sure how else to handle this error for now.
		panic(err)
	}
	reporter.HadError = true
}

func (reporter *StateErrorReporter) RuntimeError(e RuntimeError) {
	_, err := fmt.Fprintf(os.Stderr, "%s\n[line %d]\n", e, e.Token.Line)
	if err != nil { // Not sure how else to handle this error for now.
		panic(err)
	}
	reporter.HadRuntimeError = true
}
