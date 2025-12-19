// Package location tracks source code positions for tokens and AST nodes.
package location

// Position represents a token position.
type Position struct {
	Line   uint
	Column uint
	Offset uint
}

// NewPosition builds a new Position.
func NewPosition(line uint, column uint, offset uint) Position {
	return Position{
		Line:   line,
		Column: column,
		Offset: offset,
	}
}

// Location represents a token location.
type Location struct {
	File  string
	Start Position
	End   Position
}

// NewLocation build a new Location.
func NewLocation(file string, start, end Position) Location {
	return Location{
		File:  file,
		Start: start,
		End:   end,
	}
}
