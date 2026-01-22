package test

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/danielspk/tatu-lang/pkg/builder"
	"github.com/danielspk/tatu-lang/pkg/interpreter"
	"github.com/danielspk/tatu-lang/pkg/runtime"
)

const expectPrefix = "; Expect: "
const expectErrorPrefix = "; Expect Error: "

func TestPrograms(t *testing.T) {
	var files []string

	err := filepath.WalkDir("./", func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() && filepath.Ext(path) == ".tatu" {
			files = append(files, path)
		}

		return nil
	})
	if err != nil {
		t.Fatalf("exploring .tatu test files: %s", err)
	}

	for _, file := range files {
		t.Run(file, func(t *testing.T) {
			content, err := os.ReadFile(file)
			if err != nil {
				t.Errorf("reading test file: %s", err)
			}

			err = runTestSource(content, file)
			if err != nil {
				t.Errorf("running test file: %s", err)
			}
		})
	}
}

func runTestSource(source []byte, filename string) error {
	if strings.Contains(string(source), expectErrorPrefix) {
		return runErrorTest(source, filename)
	}

	return runSuccessTest(source, filename)
}

func runSuccessTest(source []byte, filename string) error {
	progBuilder := builder.NewProgramBuilder(builder.NewDefaultScanner(), builder.NewDefaultParser())
	_, ast, err := progBuilder.BuildFromFile(filename)
	if err != nil {
		return fmt.Errorf("building source: %w", err)
	}

	inter, err := interpreter.NewInterpreter()
	if err != nil {
		return fmt.Errorf("creating interpreter: %v", err)
	}

	var lastValue runtime.Value
	var checkValue string

	for _, expr := range ast.Program {
		lastValue, err = inter.Eval(expr, nil)
		if err != nil {
			return fmt.Errorf("evaluating program: %w", err)
		}
	}

	if lastValue != nil {
		checkValue = lastValue.String()

		if lastValue.Type() == runtime.StringType {
			checkValue = scapeResult(checkValue)
		}
	}

	startIdx := strings.LastIndex(string(source), expectPrefix)
	if startIdx == -1 {
		return errors.New("missing prefix value")
	}

	startIdx += 10

	endIdx := strings.LastIndex(string(source[startIdx:]), "\n")
	if endIdx == -1 {
		return errors.New("missing result value")
	}

	endIdx += startIdx

	if checkValue != string(source[startIdx:endIdx]) {
		return fmt.Errorf("expected: `%s`, found: `%s`", string(source[startIdx:endIdx]), checkValue)
	}

	return nil
}

func runErrorTest(source []byte, filename string) error {
	progBuilder := builder.NewProgramBuilder(builder.NewDefaultScanner(), builder.NewDefaultParser())
	_, ast, err := progBuilder.BuildFromFile(filename)
	if err != nil {
		return fmt.Errorf("building source: %w", err)
	}

	inter, err := interpreter.NewInterpreter()
	if err != nil {
		return fmt.Errorf("creating interpreter: %v", err)
	}

	var evalError error
	for _, expr := range ast.Program {
		_, evalError = inter.Eval(expr, nil)
		if evalError != nil {
			break
		}
	}

	if evalError == nil {
		return errors.New("expected error but program succeeded")
	}

	startIdx := strings.Index(string(source), expectErrorPrefix)
	if startIdx == -1 {
		return nil
	}

	startIdx += len(expectErrorPrefix)

	endIdx := strings.Index(string(source[startIdx:]), "\n")
	if endIdx == -1 {
		endIdx = len(source) - startIdx
	}

	expectedError := strings.TrimSpace(string(source[startIdx : startIdx+endIdx]))

	if expectedError == "" {
		return nil
	}

	actualError := evalError.Error()
	if !strings.Contains(actualError, expectedError) {
		return fmt.Errorf("expected error containing: `%s`, found: `%s`", expectedError, actualError)
	}

	return nil
}

func scapeResult(result string) string {
	result = strings.ReplaceAll(result, "\\", "\\\\")
	result = strings.ReplaceAll(result, "\n", "\\n")
	result = strings.ReplaceAll(result, "\t", "\\t")
	result = strings.ReplaceAll(result, "\r", "\\r")
	return strings.ReplaceAll(result, "\"", "\\\"")
}
