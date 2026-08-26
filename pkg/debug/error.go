// Package debug provides structured error reporting.
package debug

import (
	"fmt"

	"github.com/danielspk/tatu-lang/pkg/location"
)

// Error represents a Tatu error.
type Error struct {
	Msg string
	Loc location.Location
}

// NewError builds a new Error.
func NewError(msg string, loc location.Location) *Error {
	return &Error{
		Msg: msg,
		Loc: loc,
	}
}

// Error shows the error message.
func (e *Error) Error() string {
	if e.Loc.End.Column == 0 {
		return fmt.Sprintf("Error: %s", e.Msg)
	}

	return fmt.Sprintf("[File %s][Line %d][Column %d] Error: %s",
		e.Loc.File, e.Loc.End.Line, e.Loc.End.Column, e.Msg)
}
