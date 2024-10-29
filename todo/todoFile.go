package todo

import (
	"fmt"
	"os"
	todo "github.com/1set/todotxt"
)

// FileError represents an error that occurred during file operations
type FileError struct {
	Op   string // Operation being performed (e.g., "read", "write")
	Path string // File path where the error occurred
	Err  error  // Underlying error
}

func (e *FileError) Error() string {
	return fmt.Sprintf("todo file %s error at %s: %v", e.Op, e.Path, e.Err)
}

// ReadTodoFile reads a todo.txt file and returns a TaskList
func ReadTodoFile(path string) (todo.TaskList, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, &FileError{Op: "read", Path: path, Err: err}
	}
	defer file.Close()

	taskList, err := todo.LoadFromFile(file)
	if err != nil {
		return nil, &FileError{Op: "parse", Path: path, Err: err}
	}

	return taskList, nil
}

// WriteTodoFile writes a TaskList to a todo.txt file
func WriteTodoFile(taskList todo.TaskList, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return &FileError{Op: "create", Path: path, Err: err}
	}
	defer file.Close()

	if err := taskList.WriteToFile(file); err != nil {
		return &FileError{Op: "write", Path: path, Err: err}
	}

	return nil
}
