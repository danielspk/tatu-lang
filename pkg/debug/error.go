// Package debug provides structured error reporting.
package debug

import (
	"fmt"
	"os"
	"strings"
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

// Dump shows the error message next to reference source code.
func (e *Error) Dump() string {
	source, _ := os.ReadFile(e.File)
	lines := strings.Split(string(source), "\n")

	if len(lines) == 0 {
		return fmt.Sprintf("Error on line %d, column %d, file `%s`", e.Line, e.Column, e.File)
	}

	rawLine := ""

	if e.Line > 0 && int(e.Line) <= len(lines) {
		rawLine = lines[e.Line-1]
	}

	errColumn := int(e.Column) - 2
	if errColumn < 0 {
		errColumn = 0
	}

	errLine := strings.Repeat(" ", errColumn) + "↑" + "\n"
	errLine += strings.Repeat(" ", errColumn) + "└─ " + e.Msg

	return fmt.Sprintf("Error on line %d, column %d, file `%s`:\n\n%s\n%s", e.Line, e.Column, e.File, rawLine, errLine)
}
