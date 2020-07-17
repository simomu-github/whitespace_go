package whitespace_go

import (
	"errors"
	"fmt"
)

func parseError(parser *Parser, message string) error {
	errorMessage := fmt.Sprintf("Parse error: %s at %s:%d:%d", message, parser.filename, parser.currentLine, parser.currentColumn)
	return errors.New(errorMessage)
}

func runtimeError(executor *Executor, message string) error {
	errorMessage := fmt.Sprintf("Runtime error: %s", message)
	return errors.New(errorMessage)
}
