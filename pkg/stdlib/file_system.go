package stdlib

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/danielspk/tatu-lang/pkg/runtime"
)

// RegisterFileSystem registers file system core functions in the environment.
func RegisterFileSystem(env *runtime.Environment) error {
	functions := map[string]runtime.CoreFunction{
		"fs:read":       runtime.NewCoreFunction(fsRead),
		"fs:read-lines": runtime.NewCoreFunction(fsReadLines),
		"fs:write":      runtime.NewCoreFunction(fsWrite),
		"fs:append":     runtime.NewCoreFunction(fsAppend),
		"fs:exists":     runtime.NewCoreFunction(fsExists),
		"fs:list":       runtime.NewCoreFunction(fsList),
		"fs:mkdir":      runtime.NewCoreFunction(fsMkdir),
		"fs:move":       runtime.NewCoreFunction(fsMove),
		"fs:delete":     runtime.NewCoreFunction(fsDelete),
		"fs:is-dir":     runtime.NewCoreFunction(fsIsDir),
		"fs:size":       runtime.NewCoreFunction(fsSize),
		"fs:basename":   runtime.NewCoreFunction(fsBasename),
		"fs:temp-dir":   runtime.NewCoreFunction(fsTempDir),
	}

	for name, fn := range functions {
		if _, err := env.Define(name, fn); err != nil {
			return fmt.Errorf("failed to register file system function `%s`: %w", name, err)
		}
	}

	return nil
}

// fsRead implements the file reading core function.
// Usage: (fs:read "file.txt") => "content"
func fsRead(args ...runtime.Value) (runtime.Value, error) {
	const name = "fs:read"

	if err := expectArgs(name, 1, args); err != nil {
		return nil, err
	}

	path, err := expectString(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	content, err := os.ReadFile(path.Value)
	if err != nil {
		return nil, fmt.Errorf("`%s` failed to read file: %w", name, err)
	}

	return runtime.NewString(string(content)), nil
}

// fsReadLines implements the file reading by lines core function.
// Usage: (fs:read-lines "file.txt") => (vector "line1" "line2")
func fsReadLines(args ...runtime.Value) (runtime.Value, error) {
	const name = "fs:read-lines"

	if err := expectArgs(name, 1, args); err != nil {
		return nil, err
	}

	path, err := expectString(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	content, err := os.ReadFile(path.Value)
	if err != nil {
		return nil, fmt.Errorf("`%s` failed to read file: %w", name, err)
	}

	lines := strings.Split(string(content), "\n")
	elements := make([]runtime.Value, len(lines))

	for i, line := range lines {
		elements[i] = runtime.NewString(line)
	}

	return runtime.NewVector(elements), nil
}

// fsWrite implements the file writing core function.
// Usage: (fs:write "file.txt" "content") => nil
func fsWrite(args ...runtime.Value) (runtime.Value, error) {
	const name = "fs:write"

	if err := expectArgs(name, 2, args); err != nil {
		return nil, err
	}

	path, err := expectString(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	content, err := expectString(name, 1, args[1])
	if err != nil {
		return nil, err
	}

	if err = os.WriteFile(path.Value, []byte(content.Value), 0644); err != nil {
		return nil, fmt.Errorf("`%s` failed to write file: %w", name, err)
	}

	return runtime.NewNil(), nil
}

// fsAppend implements the file appending core function.
// Usage: (fs:append "file.txt" "more content") => nil
func fsAppend(args ...runtime.Value) (runtime.Value, error) {
	const name = "fs:append"

	if err := expectArgs(name, 2, args); err != nil {
		return nil, err
	}

	path, err := expectString(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	content, err := expectString(name, 1, args[1])
	if err != nil {
		return nil, err
	}

	file, err := os.OpenFile(path.Value, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("`%s` failed to open file: %w", name, err)
	}
	defer file.Close()

	if _, err = file.WriteString(content.Value); err != nil {
		return nil, fmt.Errorf("`%s` failed to append to file: %w", name, err)
	}

	return runtime.NewNil(), nil
}

// fsExists implements the file existence check core function.
// Usage: (fs:exists "file.txt") => true
func fsExists(args ...runtime.Value) (runtime.Value, error) {
	const name = "fs:exists"

	if err := expectArgs(name, 1, args); err != nil {
		return nil, err
	}

	path, err := expectString(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	_, err = os.Stat(path.Value)
	if err == nil {
		return runtime.NewBool(true), nil
	}
	if os.IsNotExist(err) {
		return runtime.NewBool(false), nil
	}

	return nil, fmt.Errorf("`%s` failed to check file: %w", name, err)
}

// fsList implements the directory listing core function.
// Usage: (fs:list "dir") => (vector "file1.txt" "file2.txt")
func fsList(args ...runtime.Value) (runtime.Value, error) {
	const name = "fs:list"

	if err := expectArgs(name, 1, args); err != nil {
		return nil, err
	}

	path, err := expectString(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(path.Value)
	if err != nil {
		return nil, fmt.Errorf("`%s` failed to list directory: %w", name, err)
	}

	elements := make([]runtime.Value, len(entries))

	for i, entry := range entries {
		elements[i] = runtime.NewString(entry.Name())
	}

	return runtime.NewVector(elements), nil
}

// fsMkdir implements the directory creation core function.
// Usage: (fs:mkdir "newdir") => nil
func fsMkdir(args ...runtime.Value) (runtime.Value, error) {
	const name = "fs:mkdir"

	if err := expectArgs(name, 1, args); err != nil {
		return nil, err
	}

	path, err := expectString(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	if err = os.MkdirAll(path.Value, 0755); err != nil {
		return nil, fmt.Errorf("`%s` failed to create directory: %w", name, err)
	}

	return runtime.NewNil(), nil
}

// fsMove implements the file/directory moving core function.
// Usage: (fs:move "old.txt" "new.txt") => nil
func fsMove(args ...runtime.Value) (runtime.Value, error) {
	const name = "fs:move"

	if err := expectArgs(name, 2, args); err != nil {
		return nil, err
	}

	oldPath, err := expectString(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	newPath, err := expectString(name, 1, args[1])
	if err != nil {
		return nil, err
	}

	if err = os.Rename(oldPath.Value, newPath.Value); err != nil {
		return nil, fmt.Errorf("`%s` failed to move file: %w", name, err)
	}

	return runtime.NewNil(), nil
}

// fsDelete implements the file/directory deletion core function.
// Usage: (fs:delete "file.txt") => nil
func fsDelete(args ...runtime.Value) (runtime.Value, error) {
	const name = "fs:delete"

	if err := expectArgs(name, 1, args); err != nil {
		return nil, err
	}

	path, err := expectString(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	if err = os.RemoveAll(path.Value); err != nil {
		return nil, fmt.Errorf("`%s` failed to delete: %w", name, err)
	}

	return runtime.NewNil(), nil
}

// fsIsDir implements the directory check core function.
// Usage: (fs:is-dir "path") => true
func fsIsDir(args ...runtime.Value) (runtime.Value, error) {
	const name = "fs:is-dir"

	if err := expectArgs(name, 1, args); err != nil {
		return nil, err
	}

	path, err := expectString(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	info, err := os.Stat(path.Value)
	if err != nil {
		return nil, fmt.Errorf("`%s` failed to check path: %w", name, err)
	}

	return runtime.NewBool(info.IsDir()), nil
}

// fsSize implements the file size core function.
// Usage: (fs:size "file.txt") => 1024
func fsSize(args ...runtime.Value) (runtime.Value, error) {
	const name = "fs:size"

	if err := expectArgs(name, 1, args); err != nil {
		return nil, err
	}

	path, err := expectString(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	info, err := os.Stat(path.Value)
	if err != nil {
		return nil, fmt.Errorf("`%s` failed to get file info: %w", name, err)
	}

	return runtime.NewNumber(float64(info.Size())), nil
}

// fsBasename implements the basename extraction core function.
// Usage: (fs:basename "/path/to/file.txt") => "file.txt"
func fsBasename(args ...runtime.Value) (runtime.Value, error) {
	const name = "fs:basename"

	if err := expectArgs(name, 1, args); err != nil {
		return nil, err
	}

	path, err := expectString(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	basename := filepath.Base(path.Value)

	return runtime.NewString(basename), nil
}

// fsTempDir implements the temporary directory core function.
// Usage: (fs:temp-dir) => "/tmp" or "C:\Users\...\Temp"
func fsTempDir(args ...runtime.Value) (runtime.Value, error) {
	const name = "fs:temp-dir"

	if err := expectArgs(name, 0, args); err != nil {
		return nil, err
	}

	return runtime.NewString(os.TempDir()), nil
}
