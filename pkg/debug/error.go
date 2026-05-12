// Package debug provides structured error reporting.
package debug

import (
	"fmt"
)

// Error represent a Tatu error.
type Error struct {
	Msg    string
	Line   uint
	Column uint
	File   string
}

// Error shows the error message.
func (e *Error) Error() string {
	return fmt.Sprintf("[Line %d][Column %d] Error: %s", e.Line, e.Column, e.Msg)
}
