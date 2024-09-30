package todo

import (
	"log"
	todo "github.com/1set/todotxt"
)

func readTodoFile()  {
	if _, err := todo.LoadFromPath("todo.txt"); err != nil {
		log.Fatal(err)
	}
}
