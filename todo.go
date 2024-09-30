package main

import (
	"log"
	todo "github.com/1set/todotxt"
)

func main()  {
	if _, err := todo.LoadFromPath("todo.txt"); err != nil {
		log.Fatal(err)
	}
}
