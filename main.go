package main

import (
	"os"

	"github.com/icholy/monkey/repl"
)

func main() {
	repl.Run(os.Stdin, os.Stdout)
}
