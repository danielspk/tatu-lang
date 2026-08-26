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
	defer func() {
		if r := recover(); r != nil {
			exitWithError(fmt.Errorf("unexpected failure: %v", r), nil)
		}
	}()

	printTokens := flag.Bool("printTokens", false, "print the generated tokens")
	printAST := flag.Bool("printAST", false, "print the generated AST")
	printInfo := flag.Bool("printInfo", true, "print the tatu header info")
	flag.Parse()

	if flag.NArg() == 0 {
		exitWithError(fmt.Errorf("usage `tatu [arguments] <source file>`"), nil)
	}

	filename := flag.Arg(0)

	// building from a source file
	progBuilder := builder.NewProgramBuilderWithDefaults()
	tokens, ast, err := progBuilder.BuildFromFile(filename)
	if err != nil {
		exitWithError(err, progBuilder.Sources())
	}

	if *printTokens {
		fmt.Print(pretty.DumpTokensHeader())
		fmt.Println(pretty.DumpTokens(tokens))
	}

	if *printAST {
		fmt.Print(pretty.DumpASTHeader())
		fmt.Println(pretty.DumpAST(ast))
	}

	if *printInfo {
		fmt.Println(pretty.DumpExecution(version, filename))
	}

	// evaluating by interpreter
	inter := interpreter.NewInterpreter()

	result, err := inter.EvalProgram(ast, nil)
	if err != nil {
		exitWithError(err, progBuilder.Sources())
	}

	fmt.Println(result)
}

func exitWithError(err error, sources map[string][]byte) {
	fmt.Print(pretty.DumpError(err, sources))
	os.Exit(1)
}
