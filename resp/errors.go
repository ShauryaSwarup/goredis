package resp

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidSyntax  = errors.New("Invalid Syntax")
	ErrUnexpectedType = errors.New("Unexpected Type")
)

type ParsingError struct {
	LineNum  int
	RawInput string
	Err      error
}

func (e *ParsingError) Error() string {
	return fmt.Sprintf("resp: parsing error at line %d (%q): %v",
		e.LineNum, e.RawInput, e.Err)
}

func (e *ParsingError) Unwrap() error {
	return e.Err
}
