package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/danielspk/tatu-lang/pkg/builder"
	"github.com/danielspk/tatu-lang/pkg/interpreter"
	"github.com/danielspk/tatu-lang/pkg/pretty"
)

var version = "dev-mode"

func main() {
	printTokens := flag.Bool("printTokens", false, "print the generated tokens")
	printAST := flag.Bool("printAST", false, "print the generated AST")
	printBytecode := flag.Bool("printBytecode", false, "print the byte codes")
	printInfo := flag.Bool("printInfo", true, "print the tatu header info")
	flag.Parse()

	if len(os.Args) <= 1 {
		exitWithError(fmt.Errorf("usage `tatu [arguments] <source file>`"))
	}

	filename := os.Args[len(os.Args)-1]

	// building from a source file
	progBuilder := builder.NewProgramBuilder(builder.NewDefaultScanner(), builder.NewDefaultParser())
	tokens, ast, err := progBuilder.BuildFromFile(filename)
	if err != nil {
		exitWithError(err)
	}

	// compiling to bytecode
	/*compiler := vm.NewCompiler()
	code := compiler.Compile(ast)*/

	if *printTokens {
		for _, token := range tokens {
			fmt.Println(pretty.FormatToken(token))
		}
		fmt.Println()
	}

	if *printAST {
		fmt.Println(pretty.FormatAST(ast))
	}

	if *printBytecode {
		/*for _, b := range code.Code {
			fmt.Printf("0x%02X ", b)
		}
		fmt.Println()*/
	}

	if *printInfo {
		fmt.Println(pretty.FormatRunningExecution(version, filename))
		fmt.Println(pretty.FormatRunningOutput())
	}

	// evaluating by interpreter
	inter, err := interpreter.NewInterpreter()
	if err != nil {
		exitWithError(err)
	}

	for _, expr := range ast.Program {
		result, err := inter.Eval(expr, nil)
		if err != nil {
			exitWithError(err)
		}

		fmt.Println(result)
	}

	// evaluating by virtual machine
	/*machine := vm.NewVirtualMachine()
	result, err := machine.Execute(code)
	if err != nil {
		exitWithError(err)
	}

	fmt.Println(result)*/
}

func exitWithError(err error) {
	fmt.Print(pretty.FormatError(err))
	os.Exit(1)
}
