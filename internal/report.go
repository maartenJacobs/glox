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
}

// StateErrorReporter is an implementation of ErrorReporter that tracks whether an
// error was reported and prints errors to standard error.
type StateErrorReporter struct {
	HadError bool // Whether an error has been reported.
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
