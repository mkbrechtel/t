package todo

import (
	"fmt"
	todo "github.com/1set/todotxt"
)

// PrintTaskList prints tasks with their IDs
func PrintTaskList(taskList todo.TaskList) {
	for _, t := range taskList {
		fmt.Printf("%3d %s\n", t.ID, t.String())
	}
}
