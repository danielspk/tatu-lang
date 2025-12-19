package vm

import (
	"github.com/danielspk/tatu-lang/pkg/ast"
	"github.com/danielspk/tatu-lang/pkg/runtime"
)

// Code ...
type Code struct {
	Constants []runtime.Value
	Code      []byte
}

// NewCode ...
func NewCode(name string) Code {
	return Code{
		Constants: make([]runtime.Value, 0),
		Code:      make([]byte, 0),
	}
}

// Compiler ...
type Compiler struct {
	code *Code
}

// NewCompiler ...
func NewCompiler() Compiler {
	return Compiler{
		code: &Code{},
	}
}

// Compile ...
func (c *Compiler) Compile(ast *ast.AST) *Code {
	for _, expr := range ast.Program {
		c.generate(expr)
	}

	c.code.Code = append(c.code.Code, byte(OpHalt))

	return c.code
}

// generate ...
func (c *Compiler) generate(expr ast.SExpr) {
	switch expr.Kind() {
	case ast.NumberKind:
		c.code.Code = append(c.code.Code, byte(OpConst))
		c.code.Code = append(c.code.Code, c.addNumberConstant(expr.(*ast.NumberExpr).Number))

	case ast.StringKind:
		c.code.Code = append(c.code.Code, byte(OpConst))
		c.code.Code = append(c.code.Code, c.addStringConstant(expr.(*ast.StringExpr).String))

	case ast.BoolKind, ast.SymbolKind, ast.NilKind:

	case ast.ListKind:
		if len(expr.(*ast.ListExpr).List) == 0 {
			// TODO nil support
			return
		}

		switch expr.(*ast.ListExpr).List[0].Kind() {
		case ast.SymbolKind:
			// TODO ...
		default:
			// TODO ...
		}

	default:
		// TODO ...
		//return nil, i.error("unknown expression", expr.End.Line, expr.End.Column)
	}
}

// addNumberConstant adds a number constant and returns the new constant index
func (c *Compiler) addNumberConstant(value float64) byte {
	for idx, constant := range c.code.Constants {
		if constant.Type() != runtime.NumberType {
			continue
		}

		if constant.(runtime.Number).Value == value {
			return byte(idx)
		}
	}

	// TODO overflow?

	c.code.Constants = append(c.code.Constants, runtime.NewNumber(value))

	return byte(len(c.code.Constants) - 1)
}

// addStringConstant ...
func (c *Compiler) addStringConstant(value string) byte {
	for idx, constant := range c.code.Constants {
		if constant.Type() != runtime.StringType {
			continue
		}

		if constant.(runtime.String).Value == value {
			return byte(idx)
		}
	}

	// TODO overflow?

	c.code.Constants = append(c.code.Constants, runtime.NewString(value))

	return byte(len(c.code.Constants) - 1)
}
