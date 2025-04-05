package main

import (
	"fmt"
	"os"
	
	"t5.mkbrechtel.dev/cli"
)

func main() {
	if err := cli.Execute(os.Args, os.Stdin, os.Stdout, os.Stderr); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
