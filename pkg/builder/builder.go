// Package builder orchestrates the scanning and parsing process, handling file inclusion resolution.
package builder

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/danielspk/tatu-lang/pkg/ast"
	"github.com/danielspk/tatu-lang/pkg/token"
)

// Scanner represents a lexical analyzer interface.
type Scanner interface {
	Scan(source []byte, filename string) ([]token.Token, error)
}

// Parser represents a parser interface.
type Parser interface {
	Parse(tokens []token.Token) (*ast.AST, error)
}

// ProgramBuilder is responsible for generating an AST of the program and resolving the inclusion of files and modules.
type ProgramBuilder struct {
	scanner     Scanner
	parser      Parser
	parsedFiles []string
}

// NewProgramBuilder builds a new ProgramBuilder.
func NewProgramBuilder(scanner Scanner, parser Parser) *ProgramBuilder {
	return &ProgramBuilder{
		scanner: scanner,
		parser:  parser,
	}
}

// BuildFromFile builds an AST from a file path.
func (pb *ProgramBuilder) BuildFromFile(filename string) ([]token.Token, *ast.AST, error) {
	filename = pb.fullPath(filename)

	source, err := os.ReadFile(filename)
	if err != nil {
		return nil, nil, fmt.Errorf("missing file `%s`: %w", filename, err)
	}

	return pb.BuildFromSource(source, filename)
}

// BuildFromSource builds an AST from a source code.
func (pb *ProgramBuilder) BuildFromSource(source []byte, filename string) ([]token.Token, *ast.AST, error) {
	filename = pb.fullPath(filename)

	pb.addParsedFile(filename)

	tokens, err := pb.scanner.Scan(source, filename)
	if err != nil {
		return nil, nil, fmt.Errorf("scanning source on file `%s`: %w", filename, err)
	}

	astNodes, err := pb.parser.Parse(tokens)
	if err != nil {
		return nil, nil, fmt.Errorf("parsing tokens on file `%s`: %w", filename, err)
	}

	// each ast top level node to resolve includes
	idx := 0

	for idx < len(astNodes.Program) {
		expr := astNodes.Program[idx]

		if includeFile, ok := pb.isIncludeExpr(expr); ok {
			includeFilename := pb.resolveRefPath(filename, includeFile)

			if pb.fileWasParsed(includeFilename) {
				astNodes.Program = append(astNodes.Program[:idx], astNodes.Program[idx+1:]...)

				continue
			}

			incTokens, incASTNodes, err := pb.BuildFromFile(includeFilename)
			if err != nil {
				return nil, nil, fmt.Errorf("including file `%s` in `%s`: %w", includeFilename, filename, err)
			}

			tokens = append(tokens, incTokens...)
			astNodes.Program = append(astNodes.Program[:idx], append(incASTNodes.Program, astNodes.Program[idx+1:]...)...)

			idx += len(incASTNodes.Program)
		} else {
			idx++
		}
	}

	return tokens, astNodes, nil
}

// isIncludeExpr checks if the expression is an "include" and returns the file name.
func (pb *ProgramBuilder) isIncludeExpr(expr ast.SExpr) (filename string, ok bool) {
	if expr.Kind() == ast.ListKind && len(expr.(*ast.ListExpr).List) == 2 &&
		expr.(*ast.ListExpr).List[0].Kind() == ast.SymbolKind && expr.(*ast.ListExpr).List[0].(*ast.SymbolExpr).Symbol == "include" &&
		expr.(*ast.ListExpr).List[1].Kind() == ast.StringKind {

		return expr.(*ast.ListExpr).List[1].(*ast.StringExpr).String, true
	}

	return "", false
}

// addParsedFile adds a file to the list of building files.
func (pb *ProgramBuilder) addParsedFile(filename string) {
	pb.parsedFiles = append(pb.parsedFiles, filename)
}

// fileWasParsed checks if a file was already built.
func (pb *ProgramBuilder) fileWasParsed(filename string) bool {
	for _, buildFile := range pb.parsedFiles {
		if buildFile == filename {
			return true
		}
	}

	return false
}

// fullPath resolves the absolute path of a file.
func (pb *ProgramBuilder) fullPath(filename string) string {
	if filepath.IsAbs(filename) {
		return filepath.Clean(filename)
	}

	absPath, _ := filepath.Abs(filename)

	return filepath.Clean(absPath)
}

// resolveRefPath resolves the absolute path of a destination file based on the reference file.
func (pb *ProgramBuilder) resolveRefPath(referenceFile, destinationFile string) string {
	if filepath.IsAbs(destinationFile) {
		return filepath.Clean(destinationFile)
	}

	referenceDir := filepath.Dir(referenceFile)
	destinationFile = filepath.Join(referenceDir, destinationFile)
	destinationFile, _ = filepath.Abs(destinationFile)

	return filepath.Clean(destinationFile)
}
